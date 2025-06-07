package models

import (
	"time"

	"github.com/saufiroja/fin-ai/internal/constants"
)

type Category struct {
	CategoryId    string                 `json:"category_id"`
	Name          string                 `json:"name"`
	NameEmbedding any                    `json:"name_embedding"` // type data vector for name embedding
	Description   string                 `json:"description"`
	Type          constants.TypeCategory `json:"type"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}
