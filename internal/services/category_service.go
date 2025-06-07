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

func (c *categoryService) FindCategoryById(categoryId string) (*models.Category, error) {
	c.logging.LogInfo(fmt.Sprintf("Fetching category by ID: %s", categoryId))

	category, err := c.categoryRepository.FindCategoryById(categoryId)
	if err != nil {
		c.logging.LogError(fmt.Sprintf("Failed to fetch category by ID %s: %s", categoryId, err.Error()))
		return nil, fmt.Errorf("failed to fetch category by ID %s: %w", categoryId, err)
	}

	if category == nil {
		c.logging.LogInfo(fmt.Sprintf("Category with ID %s not found", categoryId))
		return nil, nil // Return nil if not found
	}

	c.logging.LogInfo(fmt.Sprintf("Category with ID %s fetched successfully", categoryId))
	return category, nil
}

func (c *categoryService) UpdateCategoryById(categoryId string, req *requests.UpdateCategoryRequest) error {
	c.logging.LogInfo(fmt.Sprintf("Updating category with ID: %s", categoryId))

	// Fetch existing category
	existingCategory, err := c.FindCategoryById(categoryId)
	if err != nil {
		c.logging.LogError(fmt.Sprintf("Failed to fetch category by ID %s: %s", categoryId, err.Error()))
		return fmt.Errorf("failed to fetch category by ID %s: %w", categoryId, err)
	}

	if existingCategory == nil {
		c.logging.LogInfo(fmt.Sprintf("Category with ID %s not found for update", categoryId))
		return fmt.Errorf("category with ID %s not found", categoryId)
	}

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

		c.logging.LogInfo("starting to create embedding for updated category name")
		embedding := c.openaiClient.CreateEmbedding(context.Background(), input)

		if embedding != nil && embedding.Embeddings != "" {
			embeddingChan <- embedding.Embeddings
		} else {
			errorChan <- fmt.Errorf("failed to create embedding for updated category name")
		}
	}()

	// Prepare other data concurrently
	var wg sync.WaitGroup
	var timestamp time.Time

	wg.Add(1)
	go func() {
		defer wg.Done()
		timestamp = time.Now()
	}()

	// Wait for other preparations
	wg.Wait()

	// Wait for embedding result
	var embeddingResult string
	select {
	case emb := <-embeddingChan:
		embeddingResult = emb
	case err := <-errorChan:
		c.logging.LogError(fmt.Sprintf("Failed to create embedding: %s", err.Error()))
		return fmt.Errorf("failed to create embedding: %w", err)
	}

	// Update category data
	existingCategory.Name = req.Name
	existingCategory.NameEmbedding = embeddingResult
	existingCategory.Type = req.Type
	existingCategory.UpdatedAt = timestamp

	err = c.categoryRepository.UpdateCategoryById(existingCategory)
	if err != nil {
		c.logging.LogError(fmt.Sprintf("Failed to update category: %s", err.Error()))
		return fmt.Errorf("failed to update category: %w", err)
	}

	c.logging.LogInfo(fmt.Sprintf("Category with ID %s updated successfully", categoryId))
	return nil
}

func (c *categoryService) DeleteCategoryById(categoryId string) error {
	c.logging.LogInfo(fmt.Sprintf("Deleting category with ID: %s", categoryId))

	// Check if category exists
	_, err := c.FindCategoryById(categoryId)
	if err != nil {
		c.logging.LogError(fmt.Sprintf("Failed to fetch category by ID %s: %s", categoryId, err.Error()))
		return fmt.Errorf("failed to fetch category by ID %s: %w", categoryId, err)
	}

	err = c.categoryRepository.DeleteCategoryById(categoryId)
	if err != nil {
		c.logging.LogError(fmt.Sprintf("Failed to delete category by ID %s: %s", categoryId, err.Error()))
		return fmt.Errorf("failed to delete category by ID %s: %w", categoryId, err)
	}

	c.logging.LogInfo(fmt.Sprintf("Category with ID %s deleted successfully", categoryId))
	return nil
}

// UpdateCategoriesBatch updates multiple categories concurrently using goroutines
func (c *categoryService) UpdateCategoriesBatch(updates map[string]*requests.UpdateCategoryRequest) error {
	if len(updates) == 0 {
		return nil
	}

	c.logging.LogInfo(fmt.Sprintf("Updating %d categories concurrently", len(updates)))

	// Use buffered channels to prevent goroutine blocking
	resultChan := make(chan error, len(updates))

	// Update categories concurrently
	for categoryId, req := range updates {
		go func(id string, updateReq *requests.UpdateCategoryRequest) {
			err := c.UpdateCategoryById(id, updateReq)
			resultChan <- err
		}(categoryId, req)
	}

	// Collect results from all goroutines
	var errors []error
	for i := 0; i < len(updates); i++ {
		if err := <-resultChan; err != nil {
			errors = append(errors, err)
		}
	}

	// Return combined errors if any
	if len(errors) > 0 {
		return fmt.Errorf("failed to update %d out of %d categories: %v", len(errors), len(updates), errors)
	}

	c.logging.LogInfo(fmt.Sprintf("Successfully updated %d categories", len(updates)))
	return nil
}
