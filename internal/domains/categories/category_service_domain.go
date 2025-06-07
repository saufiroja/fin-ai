package categories

import (
	"github.com/saufiroja/fin-ai/internal/contracts/requests"
	"github.com/saufiroja/fin-ai/internal/contracts/responses"
	"github.com/saufiroja/fin-ai/internal/models"
)

type CategoryManager interface {
	CreateCategory(category *requests.CategoryRequest) error
	FindAllCategories(req *requests.GetAllCategoryQuery) (responses.GetAllCategoryResponse, error)
	FindCategoryById(categoryId string) (*models.Category, error)
	UpdateCategoryById(categoryId string, req *requests.UpdateCategoryRequest) error
	DeleteCategoryById(categoryId string) error
}
