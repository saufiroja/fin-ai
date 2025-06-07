package categories

import "github.com/saufiroja/fin-ai/internal/models"

type CategoryStorer interface {
	InsertCategory(category *models.Category) error
}
