package requests

import (
	"github.com/saufiroja/fin-ai/internal/constants"
)

type CategoryRequest struct {
	Name string                 `json:"name"`
	Type constants.TypeCategory `json:"type"`
}
