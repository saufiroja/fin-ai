package categories

import (
	"github.com/saufiroja/fin-ai/internal/contracts/requests"
	"github.com/saufiroja/fin-ai/internal/contracts/responses"
)

type CategoryManager interface {
	CreateCategory(category *requests.CategoryRequest) error
	FindAllCategories(req *requests.GetAllCategoryQuery) (responses.GetAllCategoryResponse, error)
}
