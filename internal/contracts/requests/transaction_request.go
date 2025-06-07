package requests

import (
	"github.com/saufiroja/fin-ai/internal/constants"
)

type TransactionRequest struct {
	TransactionId     string                 `json:"transaction_id"`
	UserId            string                 `json:"user_id"`
	CategoryId        string                 `json:"category_id"`
	Type              constants.TypeCategory `json:"type"`
	Description       string                 `json:"description"`
	Amount            int64                  `json:"amount"`
	Source            string                 `json:"source"`
	IsAutoCategorized bool                   `json:"is_auto_categorized"`
}

type GetAllTransactionsQuery struct {
	Limit    int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Offset   int    `query:"offset" validate:"omitempty,min=0"`
	Category string `query:"category" validate:"omitempty"`
	Search   string `query:"search" validate:"omitempty"`
}
