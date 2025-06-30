package requests

import (
	"time"

	"github.com/saufiroja/fin-ai/internal/constants"
)

type TransactionRequest struct {
	TransactionId        string                 `json:"-"`
	UserId               string                 `json:"user_id"`
	CategoryId           string                 `json:"category_id"`
	Type                 constants.TypeCategory `json:"type"`
	Description          string                 `json:"description"`
	Amount               int64                  `json:"amount"`
	Source               string                 `json:"source"`
	IsAutoCategorized    bool                   `json:"is_auto_categorized"`
	TransactionDate      time.Time              `json:"-"`
	AiCategoryConfidence float64                `json:"-"`
	CreatedAt            time.Time              `json:"-"`
	UpdatedAt            time.Time              `json:"-"`
	Confirmed            bool                   `json:"confirmed"`
	Discount             int64                  `json:"discount" validate:"omitempty,min=0"`
}

type UpdateTransactionRequest struct {
	UserId               string                 `json:"user_id" validate:"omitempty"`
	CategoryId           string                 `json:"category_id" validate:"omitempty"`
	Type                 constants.TypeCategory `json:"type" validate:"omitempty,oneof=income expense"`
	Description          string                 `json:"description" validate:"omitempty,max=255"`
	DescriptionEmbedding any                    `json:"description_embedding" validate:"omitempty,max=255"`
	Amount               int64                  `json:"amount" validate:"omitempty,min=0"`
	Source               string                 `json:"source" validate:"omitempty,max=255"`
	IsAutoCategorized    bool                   `json:"is_auto_categorized" validate:"omitempty"`
	AiCategoryConfidence float64                `json:"ai_category_confidence" validate:"omitempty,min=0,max=1"`
}

type GetAllTransactionsQuery struct {
	Limit     int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Offset    int    `query:"offset" validate:"omitempty,min=0"`
	Category  string `query:"category" validate:"omitempty"`
	Search    string `query:"search" validate:"omitempty"`
	StartDate string `query:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate   string `query:"end_date" validate:"omitempty,datetime=2006-01-02"`
}

type OverviewTransactionsQuery struct {
	StartDate string `query:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate   string `query:"end_date" validate:"omitempty,datetime=2006-01-02"`
}
