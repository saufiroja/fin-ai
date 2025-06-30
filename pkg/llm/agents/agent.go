package agents

import (
	"context"
	"fmt"

	"github.com/saufiroja/fin-ai/config"
	"github.com/saufiroja/fin-ai/internal/contracts/responses"
	"github.com/saufiroja/fin-ai/pkg/llm/tools"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
)

// Agent defines the interface for LLM agents
type Agent interface {
	Execute(ctx context.Context, message string, toolCtx *tools.ToolContext) (*responses.ResponseAI, error)
}

// BaseAgent provides common functionality for all agents
type BaseAgent struct {
	config       *config.AppConfig
	toolRegistry *tools.ToolRegistry
}

// NewBaseAgent creates a new base agent
func NewBaseAgent(config *config.AppConfig, toolRegistry *tools.ToolRegistry) *BaseAgent {
	return &BaseAgent{
		config:       config,
		toolRegistry: toolRegistry,
	}
}

// Execute runs the agent with the given message and context
func (ba *BaseAgent) Execute(ctx context.Context, message string, toolCtx *tools.ToolContext) (*responses.ResponseAI, error) {
	// Initialize LLM client
	llm, err := ba.initializeLLMClient(ctx)
	if err != nil {
		return nil, err
	}

	// Initialize message history and get initial response
	messageHistory, err := ba.getInitialResponse(ctx, llm, message)
	if err != nil {
		return nil, err
	}

	// Process tool calls
	messageHistory, err = ba.processToolCalls(messageHistory, toolCtx)
	if err != nil {
		return nil, err
	}

	// Generate final response
	return ba.generateFinalResponse(ctx, llm, messageHistory)
}

// initializeLLMClient initializes the LLM client with proper configuration
func (ba *BaseAgent) initializeLLMClient(ctx context.Context) (llms.Model, error) {
	geminiKey := ba.config.Gemini.ApiKey
	llm, err := googleai.New(ctx,
		googleai.WithAPIKey(geminiKey),
		googleai.WithDefaultModel("gemini-2.5-flash"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize LLM client: %w", err)
	}
	return llm, nil
}

// getInitialResponse gets the initial response from the LLM with the user message
func (ba *BaseAgent) getInitialResponse(ctx context.Context, llm llms.Model, message string) ([]llms.MessageContent, error) {
	messageHistory := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, message),
	}

	resp, err := llm.GenerateContent(ctx, messageHistory, llms.WithTools(ba.toolRegistry.GetAvailableTools()))
	if err != nil {
		return nil, fmt.Errorf("failed to generate initial response: %w", err)
	}

	// Translate the model's response into a MessageContent element
	respchoice := resp.Choices[0]
	assistantResponse := llms.TextParts(llms.ChatMessageTypeAI, respchoice.Content)
	for _, tc := range respchoice.ToolCalls {
		assistantResponse.Parts = append(assistantResponse.Parts, tc)
	}
	messageHistory = append(messageHistory, assistantResponse)

	return messageHistory, nil
}

// processToolCalls processes all tool calls from the LLM response
func (ba *BaseAgent) processToolCalls(messageHistory []llms.MessageContent, toolCtx *tools.ToolContext) ([]llms.MessageContent, error) {
	// Get the last message which should contain tool calls
	if len(messageHistory) == 0 {
		return messageHistory, nil
	}

	lastMessage := messageHistory[len(messageHistory)-1]

	// Process each tool call
	for _, part := range lastMessage.Parts {
		if toolCall, ok := part.(llms.ToolCall); ok {
			toolResponse, err := ba.toolRegistry.ExecuteTool(toolCall, toolCtx)
			if err != nil {
				return nil, fmt.Errorf("failed to execute tool call: %w", err)
			}
			messageHistory = append(messageHistory, toolResponse)
		}
	}

	return messageHistory, nil
}

// generateFinalResponse generates the final response after tool execution
func (ba *BaseAgent) generateFinalResponse(ctx context.Context, llm llms.Model, messageHistory []llms.MessageContent) (*responses.ResponseAI, error) {
	resp, err := llm.GenerateContent(ctx, messageHistory, llms.WithTools(ba.toolRegistry.GetAvailableTools()))
	if err != nil {
		return nil, fmt.Errorf("failed to generate final response: %w", err)
	}

	finalChoice := resp.Choices[0]

	// Extract token information from GenerationInfo
	inputTokens, outputTokens := ba.extractTokenCounts(finalChoice.GenerationInfo)

	return &responses.ResponseAI{
		Response:    finalChoice.Content,
		InputToken:  inputTokens,
		OutputToken: outputTokens,
	}, nil
}

// extractTokenCounts extracts input and output token counts from GenerationInfo
func (ba *BaseAgent) extractTokenCounts(generationInfo map[string]any) (int, int) {
	var inputTokens, outputTokens int

	if inputVal, exists := generationInfo["input_tokens"]; exists {
		if inputInt32, ok := inputVal.(int32); ok {
			inputTokens = int(inputInt32)
		}
	}

	if outputVal, exists := generationInfo["output_tokens"]; exists {
		if outputInt32, ok := outputVal.(int32); ok {
			outputTokens = int(outputInt32)
		}
	}

	return inputTokens, outputTokens
}
