package transaction

import (
	"github.com/saufiroja/fin-ai/internal/contracts/requests"
	"github.com/saufiroja/fin-ai/internal/models"
)

type TransactionStorer interface {
	InsertTransaction(transaction *models.Transaction) error
	GetTransactionByID(id string) (*models.Transaction, error)
	UpdateTransaction(transaction *models.Transaction) error
	DeleteTransaction(id string) error
	GetAllTransactions(req *requests.GetAllTransactionsQuery, userId string) ([]models.Transaction, error)
	CountAllTransactions(req *requests.GetAllTransactionsQuery, userId string) (int64, error)
}
