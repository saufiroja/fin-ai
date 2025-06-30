package services

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	"github.com/saufiroja/fin-ai/internal/constants/prompt"
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

func (t *transactionService) GetAllTransactions(req *requests.GetAllTransactionsQuery, userId string) (*responses.GetAllTransactionsResponse, error) {
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

	transactions, err := t.transactionRepository.GetAllTransactions(queryReq, userId)
	if err != nil {
		t.logging.LogError(fmt.Sprintf("Error fetching transactions: %v", err))
		return nil, err
	}

	// Use filtered count to get accurate count with the same filters
	count, err := t.transactionRepository.CountAllTransactions(req, userId)
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

func (t *transactionService) DeleteTransaction(id string) error {
	t.logging.LogInfo(fmt.Sprintf("Deleting transaction with ID: %s", id))

	_, err := t.GetDetailedTransaction(id)
	if err != nil {
		t.logging.LogError(fmt.Sprintf("Transaction not found for deletion: %s", id))
		return fmt.Errorf("transaction not found: %w", err)
	}

	err = t.transactionRepository.DeleteTransaction(id)
	if err != nil {
		t.logging.LogError(fmt.Sprintf("Error deleting transaction: %v", err))
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	t.logging.LogInfo(fmt.Sprintf("Transaction with ID %s deleted successfully", id))
	return nil
}

func (t *transactionService) GetDetailedTransaction(id string) (*models.Transaction, error) {
	t.logging.LogInfo(fmt.Sprintf("Fetching detailed transaction for ID: %s", id))
	transaction, err := t.transactionRepository.GetTransactionByID(id)
	if err != nil {
		t.logging.LogError(fmt.Sprintf("Error fetching transaction by ID %s: %v", id, err))
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}
	if transaction == nil {
		t.logging.LogWarn(fmt.Sprintf("Transaction with ID %s not found", id))
		return nil, fmt.Errorf("transaction not found")
	}

	t.logging.LogInfo(fmt.Sprintf("Successfully fetched transaction with ID: %s", id))
	return transaction, nil
}

func (t *transactionService) GetTransactionsStats() (*models.Transaction, error) {
	panic("unimplemented")
}

func (t *transactionService) InsertTransaction(req *requests.TransactionRequest) error {
	t.logging.LogInfo(fmt.Sprintf("Inserting transaction: %+v", req))

	// Use channels to communicate between goroutines
	embeddingChan := make(chan *responses.ResponseEmbedding)
	confidenceChan := make(chan float64)
	errorChan := make(chan error, 2) // Buffer for 2 possible errors

	// Start embedding creation in a goroutine
	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.logging.LogError(fmt.Sprintf("Panic in embedding goroutine: %v", r))
				errorChan <- fmt.Errorf("embedding creation failed: %v", r)
			}
		}()

		input := openai.EmbeddingNewParamsInputUnion{
			OfString: param.NewOpt(req.Description), // deskripsi transaksi sebagai input embedding
		}

		t.logging.LogInfo("Starting to create embedding for transaction description")
		embedding := t.openaiClient.CreateEmbedding(context.Background(), input)

		if embedding != nil && embedding.Embeddings != "" {
			embeddingChan <- embedding
		} else {
			errorChan <- fmt.Errorf("failed to create embedding")
		}
	}()

	// Start AI confidence calculation in a goroutine (only if auto-categorized)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.logging.LogError(fmt.Sprintf("Panic in confidence goroutine: %v", r))
				confidenceChan <- 0.0 // Default confidence on panic
			}
		}()

		if !req.IsAutoCategorized {
			confidenceChan <- 0.0
			return
		}

		messagePrompt := []openai.ChatCompletionMessageParamUnion{
			{OfSystem: &openai.ChatCompletionSystemMessageParam{
				Name: param.Opt[string]{Value: "system"},
				Content: openai.ChatCompletionSystemMessageParamContentUnion{
					OfString: param.NewOpt(prompt.TransactionConfidenceSystemPrompt),
				},
			},
			},
			{OfUser: &openai.ChatCompletionUserMessageParam{
				Name: param.Opt[string]{Value: "user"},
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: param.NewOpt(fmt.Sprintf(prompt.TransactionConfidenceUserPromptTemplate, req.CategoryId, req.Description)),
				},
			},
			},
		}

		t.logging.LogInfo("Starting to get AI confidence score")
		responseAi, err := t.openaiClient.SendChat(context.Background(), "gpt-4o-mini", messagePrompt)
		if err != nil {
			t.logging.LogError(fmt.Sprintf("Error creating chat completion for confidence: %v", err))
			confidenceChan <- 0.0 // Default confidence on error
		} else {
			// Parse the AI response to extract confidence score
			if responseStr, ok := responseAi.Response.(string); ok && len(responseStr) > 0 {
				if confidence, parseErr := t.parseConfidenceFromResponse(responseStr); parseErr == nil {
					confidenceChan <- confidence
				} else {
					t.logging.LogWarn(fmt.Sprintf("Failed to parse AI confidence response: %v", parseErr))
					confidenceChan <- 0.0
				}
			} else {
				confidenceChan <- 0.0
			}
		}
	}()

	// Prepare other transaction data concurrently
	var wg sync.WaitGroup
	var timestamp time.Time

	wg.Add(1)
	go func() {
		defer wg.Done()
		timestamp = time.Now()
	}()

	// Wait for other preparations
	wg.Wait()

	// Wait for both embedding and confidence results
	var embedding *responses.ResponseEmbedding
	var aiCategoryConfidence float64

	// Collect results from both goroutines
	for i := 0; i < 2; i++ {
		select {
		case emb := <-embeddingChan:
			embedding = emb
		case conf := <-confidenceChan:
			aiCategoryConfidence = conf
		case err := <-errorChan:
			t.logging.LogError(fmt.Sprintf("Error in concurrent operations: %v", err))
			return fmt.Errorf("failed to process transaction: %w", err)
		}
	}

	transaction := &models.Transaction{
		TransactionId:        ulid.Make().String(),
		UserId:               req.UserId,
		CategoryId:           req.CategoryId,
		Type:                 req.Type,
		Description:          req.Description,
		DescriptionEmbedding: embedding.Embeddings,
		Amount:               req.Amount,
		Source:               req.Source,
		TransactionDate:      timestamp, // Using prepared timestamp
		AiCategoryConfidence: aiCategoryConfidence,
		IsAutoCategorized:    req.IsAutoCategorized,
		CreatedAt:            timestamp,
		UpdatedAt:            timestamp,
		Confirmed:            req.Confirmed,
		Discount:             req.Discount,
	}

	err := t.transactionRepository.InsertTransaction(transaction)
	if err != nil {
		t.logging.LogError(fmt.Sprintf("Error inserting transaction: %v", err))
		return err
	}

	t.logging.LogInfo(fmt.Sprintf("Transaction inserted successfully with ID: %s, AI confidence: %.2f", transaction.TransactionId, aiCategoryConfidence))
	return nil
}

// Helper function to parse confidence from AI response
func (t *transactionService) parseConfidenceFromResponse(response string) (float64, error) {
	// Remove any whitespace and parse as float
	response = strings.TrimSpace(response)
	confidence, err := strconv.ParseFloat(response, 64)
	if err != nil {
		return 0.0, fmt.Errorf("failed to parse confidence: %v", err)
	}

	// Ensure confidence is within valid range
	if confidence < 0.0 {
		confidence = 0.0
	} else if confidence > 1.0 {
		confidence = 1.0
	}

	return confidence, nil
}

func (t *transactionService) UpdateTransaction(transactionId string, req *requests.UpdateTransactionRequest) error {
	t.logging.LogInfo(fmt.Sprintf("Updating transaction: %+v", req))

	existingTransaction, err := t.transactionRepository.GetTransactionByID(transactionId)
	if err != nil {
		t.logging.LogError(fmt.Sprintf("Error fetching transaction by ID %s: %v", transactionId, err))
		return fmt.Errorf("failed to get transaction for update: %w", err)
	}

	if existingTransaction == nil {
		t.logging.LogWarn(fmt.Sprintf("Transaction with ID %s not found for update", transactionId))
		return fmt.Errorf("transaction not found for update")
	}

	// If the description has changed, we need to re-create the embedding and AI confidence
	if existingTransaction.Description != req.Description {
		t.logging.LogInfo("Description has changed, re-creating embedding and AI confidence")
		// Create new embedding
		input := openai.EmbeddingNewParamsInputUnion{
			OfString: param.NewOpt(req.Description), // Use new description
		}
		t.logging.LogInfo("Starting to create new embedding for updated transaction description")
		embedding := t.openaiClient.CreateEmbedding(context.Background(), input)
		if embedding == nil || embedding.Embeddings == "" {
			t.logging.LogError("Failed to create new embedding for updated transaction")
			return fmt.Errorf("failed to create new embedding for updated transaction")
		}
		req.DescriptionEmbedding = embedding.Embeddings
		t.logging.LogInfo("Successfully created new embedding for updated transaction description")
		// Re-calculate AI confidence
		messagePrompt := []openai.ChatCompletionMessageParamUnion{
			{OfSystem: &openai.ChatCompletionSystemMessageParam{
				Name: param.Opt[string]{Value: "system"},
				Content: openai.ChatCompletionSystemMessageParamContentUnion{
					OfString: param.NewOpt(prompt.TransactionConfidenceSystemPrompt),
				},
			},
			},
			{OfUser: &openai.ChatCompletionUserMessageParam{
				Name: param.Opt[string]{Value: "user"},
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: param.NewOpt(fmt.Sprintf(prompt.TransactionConfidenceUserPromptTemplate, req.CategoryId, req.Description)),
				},
			},
			},
		}
		t.logging.LogInfo("Starting to get AI confidence score for updated transaction")

		responseAi, err := t.openaiClient.SendChat(context.Background(), "gpt-4o-mini", messagePrompt)
		if err != nil {
			t.logging.LogError(fmt.Sprintf("Error creating chat completion for updated transaction confidence: %v", err))
			return fmt.Errorf("failed to get AI confidence for updated transaction: %w", err)
		}
		if responseStr, ok := responseAi.Response.(string); ok && len(responseStr) > 0 {
			if confidence, parseErr := t.parseConfidenceFromResponse(responseStr); parseErr == nil {
				req.AiCategoryConfidence = confidence
			} else {
				t.logging.LogWarn(fmt.Sprintf("Failed to parse AI confidence response for updated transaction: %v", parseErr))
				req.AiCategoryConfidence = 0.0 // Default confidence on parse error
			}
		} else {
			t.logging.LogWarn("AI confidence response for updated transaction was empty or invalid")
			req.AiCategoryConfidence = 0.0 // Default confidence if response is invalid
		}
		t.logging.LogInfo(fmt.Sprintf("Successfully updated AI confidence for transaction ID %s: %.2f", transactionId, req.AiCategoryConfidence))
	} else {
		t.logging.LogInfo("Description has not changed, skipping embedding and AI confidence update")
		// If description hasn't changed, keep existing embedding and AI confidence
		req.DescriptionEmbedding = existingTransaction.DescriptionEmbedding
		req.AiCategoryConfidence = existingTransaction.AiCategoryConfidence
	}
	// Update timestamps
	transaction := &models.Transaction{
		TransactionId:        transactionId,
		UserId:               existingTransaction.UserId,
		CategoryId:           req.CategoryId,
		Type:                 req.Type,
		Description:          req.Description,
		DescriptionEmbedding: req.DescriptionEmbedding,
		Amount:               req.Amount,
		Source:               req.Source,
		IsAutoCategorized:    req.IsAutoCategorized,
		AiCategoryConfidence: req.AiCategoryConfidence,
		TransactionDate:      existingTransaction.TransactionDate, // Keep original date
		CreatedAt:            existingTransaction.CreatedAt,       // Keep original created at
		UpdatedAt:            time.Now(),                          // Update to current time
	}

	// Update the transaction in the repository
	err = t.transactionRepository.UpdateTransaction(transaction)
	if err != nil {
		t.logging.LogError(fmt.Sprintf("Error updating transaction: %v", err))
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	t.logging.LogInfo(fmt.Sprintf("Transaction with ID %s updated successfully", transaction.TransactionId))
	return nil
}
