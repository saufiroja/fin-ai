package models

type UserKnowledge struct {
	Transactions []*Transaction `json:"transactions,omitempty"`
	Receipts     []*Receipt     `json:"receipts,omitempty"`
	// Future additions can be added here
	// Budgets      []*models.Budget      `json:"budgets,omitempty"`
	// Goals        []*models.Goal        `json:"goals,omitempty"`
}

// RelevantFinancialData represents data with similarity scores for RAG
type RelevantFinancialData struct {
	Transactions []TransactionWithScore `json:"transactions,omitempty"`
	Receipts     []ReceiptWithScore     `json:"receipts,omitempty"`
}

type TransactionWithScore struct {
	Transaction *Transaction `json:"transaction"`
	Score       float64      `json:"score"`
}

type ReceiptWithScore struct {
	Receipt *Receipt `json:"receipt"`
	Score   float64  `json:"score"`
}
