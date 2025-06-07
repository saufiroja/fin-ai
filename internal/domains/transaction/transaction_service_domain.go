package transaction

import (
	"github.com/saufiroja/fin-ai/internal/contracts/requests"
	"github.com/saufiroja/fin-ai/internal/contracts/responses"
	"github.com/saufiroja/fin-ai/internal/models"
)

type TransactionManager interface {
	InsertTransaction(req *requests.TransactionRequest) error
	UpdateTransaction(transaction *models.Transaction) error
	DeleteTransaction(id string) error
	GetAllTransactions(*requests.GetAllTransactionsQuery) (*responses.GetAllTransactionsResponse, error)
	GetTransactionsStats() (*models.Transaction, error)
	GetDetailedTransaction(id string) (*models.Transaction, error)
}
