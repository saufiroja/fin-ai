package responses

import "github.com/saufiroja/fin-ai/internal/models"

type GetAllTransactionsResponse struct {
	TotalPages   int64                `json:"total_pages"`
	CurrentPage  int64                `json:"current_page"`
	Total        int64                `json:"total"`
	Transactions []models.Transaction `json:"transactions"`
}

type OverviewTransactions struct {
	TotalIncome  int64 `json:"total_income"`
	TotalExpense int64 `json:"total_expense"`
}

type OverviewTransactionsResponse struct {
	TotalIncome       string `json:"total_income"`
	TotalExpense      string `json:"total_expense"`
	TotalTransactions string `json:"total_transactions"`
}
