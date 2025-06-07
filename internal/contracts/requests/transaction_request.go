package requests

import (
	"time"

	"github.com/saufiroja/fin-ai/internal/constants"
)

type TransactionRequest struct {
	TransactionId        string                 `json:"transaction_id"`
	UserId               string                 `json:"user_id"`
	CategoryId           string                 `json:"category_id"`
	Type                 constants.TypeCategory `json:"type"`
	Description          string                 `json:"description"`
	DescriptionEmbedding any                    `json:"description_embedding"` // type data vector for description embedding
	Amount               int64                  `json:"amount"`
	Source               string                 `json:"source"`
	TransactionDate      time.Time              `json:"transaction_date"`
	AiCategoryConfidence float64                `json:"ai_category_confidence"` // confidence score for AI category prediction
	IsAutoCategorized    bool                   `json:"is_auto_categorized"`    // the transaction was auto-categorized by AI
	CreatedAt            time.Time              `json:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at"`
}

type GetAllTransactionsQuery struct {
	Limit    int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Offset   int    `query:"offset" validate:"omitempty,min=0"`
	Category string `query:"category" validate:"omitempty"`
	Search   string `query:"search" validate:"omitempty"`
}
