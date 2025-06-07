package models

import (
	"time"

	"github.com/saufiroja/fin-ai/internal/constants"
)

type AIRecommendation struct {
	RecommendationId   string                           `json:"recommendation_id"`
	UserId             string                           `json:"user_id"`
	RecommendationType constants.RecommendationType     `json:"recommendation_type"`
	Title              string                           `json:"title"`
	Content            string                           `json:"content"`           // type data jsonb for recommendation content
	ContentEmbedding   any                              `json:"content_embedding"` // type data vector for content embedding
	Priority           constants.RecommendationPriority `json:"priority"`          // Priority of the recommendation
	IsRead             bool                             `json:"is_read"`           // Indicates if the recommendation has been read
	ExpiredAt          *time.Time                       `json:"expired_at"`        // Optional expiration date for the recommendation
	CreatedAt          time.Time                        `json:"created_at"`
	UpdatedAt          time.Time                        `json:"updated_at"`
}
