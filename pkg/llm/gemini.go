package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/saufiroja/fin-ai/config"
	"github.com/saufiroja/fin-ai/internal/constants"
	"github.com/saufiroja/fin-ai/internal/contracts/requests"
	"github.com/saufiroja/fin-ai/internal/contracts/responses"
	"github.com/saufiroja/fin-ai/internal/domains/categories"
	"github.com/saufiroja/fin-ai/internal/domains/transaction"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"google.golang.org/genai"
)

type Gemini interface {
	Run(ctx context.Context, modelName string, messages []*genai.Content) (*responses.ResponseAI, error)
	RunAgent(ctx context.Context, message string, userId string) (*responses.ResponseAI, error)
	SetTransactionService(transactionService transaction.TransactionManager)
	SetCategoryService(categoryService categories.CategoryManager)
}

type GeminiClient struct {
	conf               *config.AppConfig
	transactionService transaction.TransactionManager
	categoryService    categories.CategoryManager
}

func NewGemini(conf *config.AppConfig) Gemini {
	return &GeminiClient{
		conf: conf,
	}
}

func (g *GeminiClient) Run(ctx context.Context, modelName string, messages []*genai.Content) (*responses.ResponseAI, error) {
	client, err := genai.NewClient(
		ctx,
		&genai.ClientConfig{
			APIKey: g.conf.Gemini.ApiKey,
		},
	)
	if err != nil {
		return nil, err
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		messages,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &responses.ResponseAI{
		Response:    result.Text(),
		InputToken:  int(result.UsageMetadata.PromptTokenCount),
		OutputToken: int(result.UsageMetadata.CandidatesTokenCount),
	}, nil
}

func (g *GeminiClient) RunAgent(ctx context.Context, message string, userId string) (*responses.ResponseAI, error) {
	geminiKey := g.conf.Gemini.ApiKey
	llm, err := googleai.New(ctx,
		googleai.WithAPIKey(geminiKey),
		googleai.WithDefaultModel("gemini-2.5-flash"),
	)
	if err != nil {
		return nil, err
	}

	messageHistory := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, message),
	}
	resp, err := llm.GenerateContent(ctx, messageHistory, llms.WithTools(g.GetAvailableTools()))
	if err != nil {
		return nil, err
	}

	// Translate the model's response into a MessageContent element that can be
	// added to messageHistory.
	respchoice := resp.Choices[0]
	assistantResponse := llms.TextParts(llms.ChatMessageTypeAI, respchoice.Content)
	for _, tc := range respchoice.ToolCalls {
		assistantResponse.Parts = append(assistantResponse.Parts, tc)
	}
	messageHistory = append(messageHistory, assistantResponse)

	// "Execute" tool calls by calling requested function
	for _, tc := range respchoice.ToolCalls {
		switch tc.FunctionCall.Name {
		case "insertTransaction":
			var args struct {
				CategoryId        string  `json:"categoryId"`
				Type              string  `json:"type"`
				Description       string  `json:"description"`
				Amount            float64 `json:"amount"`
				Source            string  `json:"source"`
				IsAutoCategorized bool    `json:"isAutoCategorized"`
				Confirmed         bool    `json:"confirmed"`
				Discount          float64 `json:"discount"`
			}
			if err := json.Unmarshal([]byte(tc.FunctionCall.Arguments), &args); err != nil {
				return nil, fmt.Errorf("failed to unmarshal tool arguments: %w", err)
			}

			// Actually call transaction service to insert the transaction
			var toolResponse llms.MessageContent
			if g.transactionService != nil {
				// Convert type string to TypeCategory
				var typeCategory constants.TypeCategory
				switch args.Type {
				case "income":
					typeCategory = constants.IncomeCategory
				case "expense":
					typeCategory = constants.ExpenseCategory
				default:
					typeCategory = constants.ExpenseCategory // default to expense
				}

				// Find best matching category based on description
				var categoryId string
				if args.CategoryId != "" {
					// Use provided category ID if available
					categoryId = args.CategoryId
				} else {
					// Find best matching category using category service
					bestCategoryId, err := g.findBestMatchingCategory(args.Description, typeCategory)
					if err != nil {
						toolResponse = llms.MessageContent{
							Role: llms.ChatMessageTypeTool,
							Parts: []llms.ContentPart{
								llms.ToolCallResponse{
									Name:    tc.FunctionCall.Name,
									Content: fmt.Sprintf("Failed to find suitable category: %s", err.Error()),
								},
							},
						}
						messageHistory = append(messageHistory, toolResponse)
						continue
					}
					categoryId = bestCategoryId
				}

				// Create transaction request
				transactionReq := &requests.TransactionRequest{
					UserId:            userId, // Use userId from parameter instead of args.UserId
					CategoryId:        categoryId,
					Type:              typeCategory,
					Description:       args.Description,
					Amount:            int64(args.Amount * 100), // Convert to cents
					Source:            args.Source,
					IsAutoCategorized: true, // Set to true since we're auto-categorizing
					Confirmed:         args.Confirmed,
					Discount:          int64(args.Discount * 100), // Convert to cents
				}

				// Insert transaction via service
				err = g.transactionService.InsertTransaction(transactionReq)
				if err != nil {
					toolResponse = llms.MessageContent{
						Role: llms.ChatMessageTypeTool,
						Parts: []llms.ContentPart{
							llms.ToolCallResponse{
								Name:    tc.FunctionCall.Name,
								Content: fmt.Sprintf("Failed to insert transaction: %s", err.Error()),
							},
						},
					}
				} else {
					toolResponse = llms.MessageContent{
						Role: llms.ChatMessageTypeTool,
						Parts: []llms.ContentPart{
							llms.ToolCallResponse{
								Name:    tc.FunctionCall.Name,
								Content: fmt.Sprintf("Transaction '%s' for amount %.2f successfully inserted for user %s with auto-categorized category", args.Description, args.Amount, userId),
							},
						},
					}
				}
			} else {
				toolResponse = llms.MessageContent{
					Role: llms.ChatMessageTypeTool,
					Parts: []llms.ContentPart{
						llms.ToolCallResponse{
							Name:    tc.FunctionCall.Name,
							Content: "Transaction service not available - please configure transaction service",
						},
					},
				}
			}
			messageHistory = append(messageHistory, toolResponse)

		default:
			return nil, fmt.Errorf("got unexpected function call: %v", tc.FunctionCall.Name)
		}
	}

	// Generate final response after tool execution
	resp, err = llm.GenerateContent(ctx, messageHistory, llms.WithTools(g.GetAvailableTools()))
	if err != nil {
		return nil, err
	}

	// Return the final response
	finalChoice := resp.Choices[0]
	return &responses.ResponseAI{
		Response:    finalChoice.Content,
		InputToken:  0, // TODO: Add proper token counting for agent mode
		OutputToken: 0, // TODO: Add proper token counting for agent mode
	}, nil
}

func (g *GeminiClient) GetAvailableTools() []llms.Tool {
	return []llms.Tool{
		{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "insertTransaction",
				Description: "Insert a new transaction into the system",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"categoryId": map[string]any{
							"type":        "string",
							"description": "The category ID for the transaction (optional - will be auto-detected if not provided)",
						},
						"type": map[string]any{
							"type":        "string",
							"description": "The transaction type (income/expense)",
						},
						"description": map[string]any{
							"type":        "string",
							"description": "The transaction description",
						},
						"amount": map[string]any{
							"type":        "number",
							"description": "The transaction amount",
						},
						"source": map[string]any{
							"type":        "string",
							"description": "The transaction source",
						},
						"isAutoCategorized": map[string]any{
							"type":        "boolean",
							"description": "Whether the transaction is auto-categorized",
						},
						"confirmed": map[string]any{
							"type":        "boolean",
							"description": "Whether the transaction is confirmed",
						},
						"discount": map[string]any{
							"type":        "number",
							"description": "The discount amount",
						},
					},
					"required": []string{"type", "description", "amount", "source"},
				},
			},
		},
	}
}

func (g *GeminiClient) SetTransactionService(transactionService transaction.TransactionManager) {
	g.transactionService = transactionService
}

func (g *GeminiClient) SetCategoryService(categoryService categories.CategoryManager) {
	g.categoryService = categoryService
}

// findBestMatchingCategory finds the best matching category based on transaction description
func (g *GeminiClient) findBestMatchingCategory(description string, transactionType constants.TypeCategory) (string, error) {
	if g.categoryService == nil {
		return "", fmt.Errorf("category service not available")
	}

	// Get all categories with a high limit to search through all available categories
	req := &requests.GetAllCategoryQuery{
		Offset: 1,
		Limit:  1000, // High limit to get all categories
		Search: "",   // No search filter, get all
	}

	categoriesResponse, err := g.categoryService.FindAllCategories(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch categories: %w", err)
	}

	if len(categoriesResponse.Categories) == 0 {
		return "", fmt.Errorf("no categories available")
	}

	// Filter categories by transaction type and find best match
	var bestCategoryId string
	var bestScore float64 = 0

	for _, category := range categoriesResponse.Categories {
		// Only consider categories that match the transaction type
		if category.Type == transactionType {
			// Simple keyword matching - you could enhance this with embedding similarity
			score := g.calculateCategoryScore(description, category.Name)
			if score > bestScore {
				bestScore = score
				bestCategoryId = category.CategoryId
			}
		}
	}

	if bestCategoryId == "" {
		// If no match found, return the first category of the correct type
		for _, category := range categoriesResponse.Categories {
			if category.Type == transactionType {
				bestCategoryId = category.CategoryId
				break
			}
		}
	}

	if bestCategoryId == "" {
		return "", fmt.Errorf("no suitable category found for transaction type %v", transactionType)
	}

	return bestCategoryId, nil
}

// calculateCategoryScore calculates a simple similarity score between transaction description and category name
func (g *GeminiClient) calculateCategoryScore(description, categoryName string) float64 {
	// Simple keyword matching - convert to lowercase and check for common words
	descWords := strings.Fields(strings.ToLower(description))
	catWords := strings.Fields(strings.ToLower(categoryName))

	matches := 0
	for _, descWord := range descWords {
		for _, catWord := range catWords {
			if strings.Contains(descWord, catWord) || strings.Contains(catWord, descWord) {
				matches++
			}
		}
	}

	if len(descWords) == 0 {
		return 0
	}

	return float64(matches) / float64(len(descWords))
}
