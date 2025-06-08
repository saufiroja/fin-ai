package models

import (
	"time"

	"github.com/saufiroja/fin-ai/internal/constants"
)

type Transaction struct {
	TransactionId        string                 `json:"transaction_id"`
	UserId               string                 `json:"user_id"`
	CategoryId           string                 `json:"category_id"`
	Type                 constants.TypeCategory `json:"type"`
	Description          string                 `json:"description"`
	DescriptionEmbedding any                    `json:"-"` // type data vector for description embedding
	Amount               int64                  `json:"amount"`
	Source               string                 `json:"source"`
	TransactionDate      time.Time              `json:"transaction_date"`
	AiCategoryConfidence float64                `json:"ai_category_confidence"` // confidence score for AI category prediction
	IsAutoCategorized    bool                   `json:"is_auto_categorized"`    // the transaction was auto-categorized by AI
	CreatedAt            time.Time              `json:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at"`
	Confirmed            bool                   `json:"confirmed"` // whether the transaction is confirmed by the user
}
