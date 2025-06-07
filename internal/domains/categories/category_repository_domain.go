package categories

import (
	"github.com/saufiroja/fin-ai/internal/contracts/requests"
	"github.com/saufiroja/fin-ai/internal/models"
)

type CategoryStorer interface {
	InsertCategory(*models.Category) error
	FindAllCategories(*requests.GetAllCategoryQuery) ([]models.Category, error)
	CountCategories(*requests.GetAllCategoryQuery) (int64, error)
	FindCategoryById(categoryId string) (*models.Category, error)
	UpdateCategoryById(*models.Category) error
	DeleteCategoryById(categoryId string) error
}
