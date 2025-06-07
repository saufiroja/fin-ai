package categories

import (
	"github.com/saufiroja/fin-ai/internal/contracts/requests"
)

type CategoryManager interface {
	CreateCategory(category *requests.CategoryRequest) error
}
