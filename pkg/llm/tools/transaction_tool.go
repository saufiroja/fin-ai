package tools

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/saufiroja/fin-ai/internal/constants"
	"github.com/saufiroja/fin-ai/internal/contracts/requests"
	"github.com/tmc/langchaingo/llms"
)

// TransactionTool handles transaction-related tool calls
type TransactionTool struct{}

// NewTransactionTool creates a new transaction tool
func NewTransactionTool() *TransactionTool {
	return &TransactionTool{}
}

// Name returns the tool name
func (tt *TransactionTool) Name() string {
	return "insertTransaction"
}

// Handle handles the transaction tool call
func (tt *TransactionTool) Handle(toolCall llms.ToolCall, ctx *ToolContext) (llms.MessageContent, error) {
	var args TransactionArgs
	if err := json.Unmarshal([]byte(toolCall.FunctionCall.Arguments), &args); err != nil {
		return CreateErrorToolResponse(toolCall.FunctionCall.Name, fmt.Sprintf("Failed to parse arguments: %s", err.Error())), nil
	}

	// Check if transaction service is available
	if ctx.TransactionService == nil {
		return CreateErrorToolResponse(toolCall.FunctionCall.Name, "Transaction service not available"), nil
	}

	// Convert type and find category
	typeCategory := tt.convertTransactionType(args.Type)
	categoryId, err := tt.determineCategoryId(args.CategoryId, args.Description, typeCategory, ctx)
	if err != nil {
		return CreateErrorToolResponse(toolCall.FunctionCall.Name, fmt.Sprintf("Failed to determine category: %s", err.Error())), nil
	}

	// Create and execute transaction
	transactionReq := &requests.TransactionRequest{
		UserId:            ctx.UserId,
		CategoryId:        categoryId,
		Type:              typeCategory,
		Description:       args.Description,
		Amount:            int64(args.Amount), // Convert to Rupiah integer
		Source:            args.Source,
		IsAutoCategorized: args.CategoryId == "", // Auto-categorized if no category provided
		Confirmed:         args.Confirmed,
		Discount:          int64(args.Discount), // Convert to Rupiah integer
	}

	err = ctx.TransactionService.InsertTransaction(transactionReq)
	if err != nil {
		return CreateErrorToolResponse(toolCall.FunctionCall.Name, fmt.Sprintf("Failed to insert transaction: %s", err.Error())), nil
	}

	return CreateSuccessToolResponse(toolCall.FunctionCall.Name, fmt.Sprintf("Transaction '%s' for amount %.2f successfully inserted for user %s", args.Description, args.Amount, ctx.UserId)), nil
}

// convertTransactionType converts string type to TypeCategory
func (tt *TransactionTool) convertTransactionType(typeStr string) constants.TypeCategory {
	switch typeStr {
	case "income":
		return constants.IncomeCategory
	case "expense":
		return constants.ExpenseCategory
	default:
		return constants.ExpenseCategory // default to expense
	}
}

// determineCategoryId determines the category ID to use
func (tt *TransactionTool) determineCategoryId(providedCategoryId, description string, transactionType constants.TypeCategory, ctx *ToolContext) (string, error) {
	if providedCategoryId != "" {
		return providedCategoryId, nil
	}

	// Find best matching category using category service
	return tt.findBestMatchingCategory(description, transactionType, ctx)
}

// findBestMatchingCategory finds the best matching category based on transaction description
func (tt *TransactionTool) findBestMatchingCategory(description string, transactionType constants.TypeCategory, ctx *ToolContext) (string, error) {
	if ctx.CategoryService == nil {
		return "", fmt.Errorf("category service not available")
	}

	// Get all categories with a high limit to search through all available categories
	req := &requests.GetAllCategoryQuery{
		Offset: 1,
		Limit:  1000, // High limit to get all categories
		Search: "",   // No search filter, get all
	}

	categoriesResponse, err := ctx.CategoryService.FindAllCategories(req)
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
			score := tt.calculateCategoryScore(description, category.Name)
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
func (tt *TransactionTool) calculateCategoryScore(description, categoryName string) float64 {
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
