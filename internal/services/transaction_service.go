package services

import (
	"fmt"
	"math"

	"github.com/saufiroja/fin-ai/internal/contracts/requests"
	"github.com/saufiroja/fin-ai/internal/contracts/responses"
	"github.com/saufiroja/fin-ai/internal/domains/transaction"
	"github.com/saufiroja/fin-ai/internal/models"
	logging "github.com/saufiroja/fin-ai/pkg/loggings"
)

type transactionService struct {
	transactionRepository transaction.TransactionStorer
	logging               logging.Logger
}

func NewTransactionService(
	transactionRepository transaction.TransactionStorer,
	logging logging.Logger,
) transaction.TransactionManager {
	return &transactionService{
		transactionRepository: transactionRepository,
		logging:               logging,
	}
}

func (t *transactionService) GetAllTransactions(req *requests.GetAllTransactionsQuery) (*responses.GetAllTransactionsResponse, error) {
	t.logging.LogInfo(fmt.Sprintf("Fetching all transactions with query: %+v", req))

	// Calculate offset for pagination (convert page-based to offset-based)
	offset := 0
	if req.Offset > 1 {
		offset = (req.Offset - 1) * req.Limit
	}

	// Create query with proper offset
	queryReq := &requests.GetAllTransactionsQuery{
		Offset:   offset,
		Limit:    req.Limit,
		Category: req.Category,
		Search:   req.Search,
	}

	transactions, err := t.transactionRepository.GetAllTransactions(queryReq)
	if err != nil {
		t.logging.LogError(fmt.Sprintf("Error fetching transactions: %v", err))
		return nil, err
	}

	// Use filtered count to get accurate count with the same filters
	count, err := t.transactionRepository.CountAllTransactions(req)
	if err != nil {
		t.logging.LogError(fmt.Sprintf("Error counting transactions: %v", err))
		return nil, err
	}

	totalPages := math.Max(1, math.Ceil(float64(count)/float64(req.Limit)))
	currentPage := math.Min(float64(req.Offset), float64(totalPages))

	res := &responses.GetAllTransactionsResponse{
		Transactions: transactions,
		CurrentPage:  int64(currentPage),
		TotalPages:   int64(totalPages),
		Total:        int64(count),
	}

	t.logging.LogInfo(fmt.Sprintf("Successfully fetched %d transactions", len(transactions)))
	return res, nil
}

// DeleteTransaction implements transaction.TransactionManager.
func (t *transactionService) DeleteTransaction(id string) error {
	return t.transactionRepository.DeleteTransaction(id)
}

// GetDetailedTransaction implements transaction.TransactionManager.
func (t *transactionService) GetDetailedTransaction(id string) (*models.Transaction, error) {
	return t.transactionRepository.GetTransactionByID(id)
}

// GetTransactionsStats implements transaction.TransactionManager.
func (t *transactionService) GetTransactionsStats() (*models.Transaction, error) {
	panic("unimplemented")
}

// InsertTransaction implements transaction.TransactionManager.
func (t *transactionService) InsertTransaction(transaction *models.Transaction) error {
	panic("unimplemented")
}

// UpdateTransaction implements transaction.TransactionManager.
func (t *transactionService) UpdateTransaction(transaction *models.Transaction) error {
	panic("unimplemented")
}
