package agents

import (
	"context"

	"github.com/saufiroja/fin-ai/config"
	"github.com/saufiroja/fin-ai/internal/contracts/responses"
	"github.com/saufiroja/fin-ai/pkg/llm/tools"
)

// TransactionAgent specializes in handling transaction-related tasks
type TransactionAgent struct {
	*BaseAgent
}

// NewTransactionAgent creates a new transaction agent
func NewTransactionAgent(config *config.AppConfig) *TransactionAgent {
	// Create tool registry and register transaction tools
	toolRegistry := tools.NewToolRegistry()
	toolRegistry.RegisterTool(tools.NewTransactionTool())

	baseAgent := NewBaseAgent(config, toolRegistry)
	return &TransactionAgent{
		BaseAgent: baseAgent,
	}
}

// Execute runs the transaction agent with enhanced transaction capabilities
func (ta *TransactionAgent) Execute(ctx context.Context, message string, toolCtx *tools.ToolContext) (*responses.ResponseAI, error) {
	// For now, just use the base agent functionality
	// In the future, you could add transaction-specific pre/post processing here
	return ta.BaseAgent.Execute(ctx, message, toolCtx)
}
