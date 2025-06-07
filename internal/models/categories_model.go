package models

import (
	"time"

	"github.com/saufiroja/fin-ai/internal/constants"
)

type Category struct {
	CategoryId    string                 `json:"category_id"`
	Name          string                 `json:"name"`
	NameEmbedding any                    `json:"-"` // type data vector for name embedding
	Type          constants.TypeCategory `json:"type"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}
