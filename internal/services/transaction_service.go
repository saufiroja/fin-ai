package services

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	"github.com/saufiroja/fin-ai/internal/contracts/requests"
	"github.com/saufiroja/fin-ai/internal/contracts/responses"
	"github.com/saufiroja/fin-ai/internal/domains/transaction"
	"github.com/saufiroja/fin-ai/internal/models"
	"github.com/saufiroja/fin-ai/pkg/llm"
	logging "github.com/saufiroja/fin-ai/pkg/loggings"
)

type transactionService struct {
	transactionRepository transaction.TransactionStorer
	logging               logging.Logger
	openaiClient          llm.OpenAI
}

func NewTransactionService(
	transactionRepository transaction.TransactionStorer,
	logging logging.Logger,
	openaiClient llm.OpenAI,
) transaction.TransactionManager {
	return &transactionService{
		transactionRepository: transactionRepository,
		logging:               logging,
		openaiClient:          openaiClient,
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

func (t *transactionService) InsertTransaction(req *requests.TransactionRequest) error {
	t.logging.LogInfo(fmt.Sprintf("Inserting transaction: %+v", req))

	input := openai.EmbeddingNewParamsInputUnion{
		OfString: param.NewOpt(req.Description), // deskripsi transaksi sebagai input embedding
	}
	// create embedding if not provided
	embedding := t.openaiClient.CreateEmbedding(context.Background(), input)

	transaction := &models.Transaction{
		TransactionId:        ulid.Make().String(),
		UserId:               req.UserId,
		CategoryId:           req.CategoryId,
		Type:                 req.Type,
		Description:          req.Description,
		DescriptionEmbedding: embedding.Embeddings,
		Amount:               req.Amount,
		Source:               req.Source,
		TransactionDate:      req.TransactionDate,
		AiCategoryConfidence: req.AiCategoryConfidence,
		IsAutoCategorized:    req.IsAutoCategorized,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	err := t.transactionRepository.InsertTransaction(transaction)
	if err != nil {
		t.logging.LogError(fmt.Sprintf("Error inserting transaction: %v", err))
		return err
	}

	return nil
}

// UpdateTransaction implements transaction.TransactionManager.
func (t *transactionService) UpdateTransaction(transaction *models.Transaction) error {
	panic("unimplemented")
}
