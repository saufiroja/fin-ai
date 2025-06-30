package tools

import (
	"github.com/saufiroja/fin-ai/internal/domains/categories"
	"github.com/saufiroja/fin-ai/internal/domains/transaction"
	"github.com/tmc/langchaingo/llms"
)

// ToolContext contains all dependencies needed for tool execution
type ToolContext struct {
	TransactionService transaction.TransactionManager
	CategoryService    categories.CategoryManager
	UserId             string
}

// ToolHandler defines the interface for handling tool calls
type ToolHandler interface {
	Name() string
	Handle(toolCall llms.ToolCall, ctx *ToolContext) (llms.MessageContent, error)
}

// TransactionArgs represents the arguments for transaction tool
type TransactionArgs struct {
	CategoryId        string  `json:"categoryId"`
	Type              string  `json:"type"`
	Description       string  `json:"description"`
	Amount            float64 `json:"amount"`
	Source            string  `json:"source"`
	IsAutoCategorized bool    `json:"isAutoCategorized"`
	Confirmed         bool    `json:"confirmed"`
	Discount          float64 `json:"discount"`
}
