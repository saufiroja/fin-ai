package models

import "time"

type Budget struct {
	BudgetId    string    `json:"budget_id"`
	UserId      string    `json:"user_id"`
	CategoryId  string    `json:"category_id"`
	AmountLimit int64     `json:"amount_limit"` // The maximum amount allowed for the budget
	Month       int       `json:"month"`        // The month for which the budget is set (1-12)
	Year        int       `json:"year"`         // The year for which the budget is set
	CreatedAt   time.Time `json:"created_at"`   // Timestamp when the budget was created
	UpdatedAt   time.Time `json:"updated_at"`   // Timestamp when the budget was last updated
}
