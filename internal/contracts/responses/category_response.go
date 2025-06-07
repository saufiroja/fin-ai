package responses

import "github.com/saufiroja/fin-ai/internal/models"

type GetAllCategoryResponse struct {
	TotalPages  int64             `json:"total_pages"`
	CurrentPage int64             `json:"current_page"`
	Total       int64             `json:"total"`
	Categories  []models.Category `json:"categories"`
}
