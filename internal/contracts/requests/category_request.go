package requests

import (
	"github.com/saufiroja/fin-ai/internal/constants"
)

type CategoryRequest struct {
	Name string                 `json:"name"`
	Type constants.TypeCategory `json:"type"`
}

type GetAllCategoryQuery struct {
	Limit  int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Offset int    `query:"offset" validate:"omitempty,min=0"`
	Search string `query:"search" validate:"omitempty"`
}

type UpdateCategoryRequest struct {
	Name string                 `json:"name" validate:"required"`
	Type constants.TypeCategory `json:"type" validate:"required"`
}
