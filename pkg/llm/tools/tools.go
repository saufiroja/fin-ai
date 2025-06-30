package tools

import (
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

// ToolRegistry manages all available tools
type ToolRegistry struct {
	handlers map[string]ToolHandler
}

// NewToolRegistry creates a new tool registry
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		handlers: make(map[string]ToolHandler),
	}
}

// RegisterTool registers a new tool handler
func (tr *ToolRegistry) RegisterTool(handler ToolHandler) {
	tr.handlers[handler.Name()] = handler
}

// ExecuteTool executes a tool by name
func (tr *ToolRegistry) ExecuteTool(toolCall llms.ToolCall, ctx *ToolContext) (llms.MessageContent, error) {
	handler, exists := tr.handlers[toolCall.FunctionCall.Name]
	if !exists {
		return llms.MessageContent{}, fmt.Errorf("unknown tool: %s", toolCall.FunctionCall.Name)
	}

	return handler.Handle(toolCall, ctx)
}

// GetAvailableTools returns all available tools as LLM tool definitions
func (tr *ToolRegistry) GetAvailableTools() []llms.Tool {
	var tools []llms.Tool

	// Add transaction tool
	tools = append(tools, llms.Tool{
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
	})

	return tools
}

// CreateErrorToolResponse creates an error tool response
func CreateErrorToolResponse(toolName, errorMessage string) llms.MessageContent {
	return llms.MessageContent{
		Role: llms.ChatMessageTypeTool,
		Parts: []llms.ContentPart{
			llms.ToolCallResponse{
				Name:    toolName,
				Content: errorMessage,
			},
		},
	}
}

// CreateSuccessToolResponse creates a success tool response
func CreateSuccessToolResponse(toolName, successMessage string) llms.MessageContent {
	return llms.MessageContent{
		Role: llms.ChatMessageTypeTool,
		Parts: []llms.ContentPart{
			llms.ToolCallResponse{
				Name:    toolName,
				Content: successMessage,
			},
		},
	}
}
