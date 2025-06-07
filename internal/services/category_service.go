package services

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	"github.com/saufiroja/fin-ai/internal/contracts/requests"
	"github.com/saufiroja/fin-ai/internal/contracts/responses"
	"github.com/saufiroja/fin-ai/internal/domains/categories"
	"github.com/saufiroja/fin-ai/internal/models"
	"github.com/saufiroja/fin-ai/pkg/llm"
	logging "github.com/saufiroja/fin-ai/pkg/loggings"
)

type categoryService struct {
	categoryRepository categories.CategoryStorer
	logging            logging.Logger
	openaiClient       llm.OpenAI
}

func NewCategoryService(categoryRepository categories.CategoryStorer, logging logging.Logger, openaiClient llm.OpenAI) categories.CategoryManager {
	return &categoryService{
		categoryRepository: categoryRepository,
		logging:            logging,
		openaiClient:       openaiClient,
	}
}

func (c *categoryService) CreateCategory(req *requests.CategoryRequest) error {
	c.logging.LogInfo(fmt.Sprintf("Creating category with name: %s", req.Name))

	// Use channels to communicate between goroutines
	embeddingChan := make(chan string)
	errorChan := make(chan error, 1)

	// Start embedding creation in a goroutine
	go func() {
		defer close(embeddingChan)
		defer close(errorChan)

		input := openai.EmbeddingNewParamsInputUnion{
			OfString: param.NewOpt(req.Name), // deskripsi transaksi sebagai input embedding
		}

		c.logging.LogInfo("starting to create embedding for category name")
		embedding := c.openaiClient.CreateEmbedding(context.Background(), input)

		if embedding != nil && embedding.Embeddings != "" {
			embeddingChan <- embedding.Embeddings
		} else {
			errorChan <- fmt.Errorf("failed to create embedding")
		}
	}()

	// Prepare other category data concurrently
	var wg sync.WaitGroup
	var categoryId string
	var timestamp time.Time

	wg.Add(1)
	go func() {
		defer wg.Done()
		categoryId = ulid.Make().String()
		timestamp = time.Now()
	}()

	// Wait for embedding creation and other preparations
	wg.Wait()

	// Wait for embedding result
	var embedding string
	select {
	case emb := <-embeddingChan:
		embedding = emb
	case err := <-errorChan:
		c.logging.LogError(fmt.Sprintf("Failed to create embedding: %s", err.Error()))
		return fmt.Errorf("failed to create embedding: %w", err)
	}

	newCategory := &models.Category{
		CategoryId:    categoryId,
		Name:          req.Name,
		NameEmbedding: embedding,
		Type:          req.Type,
		CreatedAt:     timestamp,
		UpdatedAt:     timestamp,
	}

	err := c.categoryRepository.InsertCategory(newCategory)
	if err != nil {
		c.logging.LogError(fmt.Sprintf("Failed to create category: %s", err.Error()))
		return fmt.Errorf("failed to create category: %w", err)
	}

	c.logging.LogInfo(fmt.Sprintf("Category created successfully with ID: %s", newCategory.CategoryId))
	return nil
}

func (c *categoryService) FindAllCategories(req *requests.GetAllCategoryQuery) (responses.GetAllCategoryResponse, error) {
	c.logging.LogInfo("Fetching all categories")
	offset := 0
	if req.Offset > 1 {
		offset = (req.Offset - 1) * req.Limit
	}

	// Create query with proper offset
	queryReq := &requests.GetAllCategoryQuery{
		Offset: offset,
		Limit:  req.Limit,
		Search: req.Search,
	}

	categoriesList, err := c.categoryRepository.FindAllCategories(queryReq)
	if err != nil {
		c.logging.LogError(fmt.Sprintf("Failed to fetch categories: %s", err.Error()))
		return responses.GetAllCategoryResponse{}, fmt.Errorf("failed to fetch categories: %w", err)
	}

	count, err := c.categoryRepository.CountCategories(queryReq)
	if err != nil {
		c.logging.LogError(fmt.Sprintf("Failed to count categories: %s", err.Error()))
		return responses.GetAllCategoryResponse{}, fmt.Errorf("failed to count categories: %w", err)
	}

	totalPages := math.Max(1, math.Ceil(float64(count)/float64(req.Limit)))
	currentPage := math.Min(float64(req.Offset), float64(totalPages))

	response := responses.GetAllCategoryResponse{
		CurrentPage: int64(currentPage),
		TotalPages:  int64(totalPages),
		Total:       int64(count),
		Categories:  categoriesList,
	}

	c.logging.LogInfo(fmt.Sprintf("Fetched %d categories successfully", len(categoriesList)))
	return response, nil
}
