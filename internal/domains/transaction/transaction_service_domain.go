package transaction

import (
	"github.com/saufiroja/fin-ai/internal/contracts/requests"
	"github.com/saufiroja/fin-ai/internal/contracts/responses"
	"github.com/saufiroja/fin-ai/internal/models"
)

type TransactionManager interface {
	InsertTransaction(req *requests.TransactionRequest) error
	UpdateTransaction(transactionId string, req *requests.UpdateTransactionRequest) error
	DeleteTransaction(id string) error
	GetAllTransactions(req *requests.GetAllTransactionsQuery, userId string) (*responses.GetAllTransactionsResponse, error)
	GetTransactionsStats() (*models.Transaction, error)
	GetDetailedTransaction(id string) (*models.Transaction, error)
	OverviewTransactions(userId string, req *requests.OverviewTransactionsQuery) (*responses.OverviewTransactionsResponse, error)
}
