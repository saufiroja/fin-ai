package categories

import (
	"github.com/saufiroja/fin-ai/internal/contracts/requests"
	"github.com/saufiroja/fin-ai/internal/models"
)

type CategoryStorer interface {
	InsertCategory(category *models.Category) error
	FindAllCategories(req *requests.GetAllCategoryQuery) ([]models.Category, error)
	CountCategories(req *requests.GetAllCategoryQuery) (int64, error)
}
