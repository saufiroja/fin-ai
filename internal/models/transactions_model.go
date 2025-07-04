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
	DescriptionEmbedding any                    `json:"-"`
	Amount               int64                  `json:"amount"`
	Source               string                 `json:"source"`
	TransactionDate      time.Time              `json:"transaction_date"`
	AiCategoryConfidence float64                `json:"ai_category_confidence"`
	IsAutoCategorized    bool                   `json:"is_auto_categorized"`
	CreatedAt            time.Time              `json:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at"`
	Confirmed            bool                   `json:"confirmed"`
	Discount             int64                  `json:"discount" validate:"omitempty,min=0"`
	PaymentMethod        string                 `json:"payment_method"`
}
