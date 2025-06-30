package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/saufiroja/fin-ai/config"
	"github.com/saufiroja/fin-ai/internal/contracts/responses"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"google.golang.org/genai"
)

type Gemini interface {
	Run(ctx context.Context, modelName string, messages []*genai.Content) (*responses.ResponseAI, error)
	RunAgent(ctx context.Context, message string) (*responses.ResponseAI, error)
}

type GeminiClient struct {
	conf *config.AppConfig
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

func (g *GeminiClient) RunAgent(ctx context.Context, message string) (*responses.ResponseAI, error) {
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
		case "getCurrentWeather":
			var args struct {
				Location string `json:"location"`
			}
			if err := json.Unmarshal([]byte(tc.FunctionCall.Arguments), &args); err != nil {
				return nil, fmt.Errorf("failed to unmarshal tool arguments: %w", err)
			}
			if strings.Contains(args.Location, "Chicago") {
				toolResponse := llms.MessageContent{
					Role: llms.ChatMessageTypeTool,
					Parts: []llms.ContentPart{
						llms.ToolCallResponse{
							Name:    tc.FunctionCall.Name,
							Content: "64 and sunny",
						},
					},
				}
				messageHistory = append(messageHistory, toolResponse)
			}
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
				Name:        "getCurrentWeather",
				Description: "Get the current weather in a given location",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"location": map[string]any{
							"type":        "string",
							"description": "The city and state, e.g. San Francisco, CA",
						},
					},
					"required": []string{"location"},
				},
			},
		},
	}
}
