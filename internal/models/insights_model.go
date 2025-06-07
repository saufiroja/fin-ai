package models

import "time"

type Insight struct {
	InsightId        string    `json:"insight_id"`        // Unique identifier for the insight
	UserId           string    `json:"user_id"`           // Identifier for the user associated with the insight
	InsightType      string    `json:"insight_type"`      // Type of insight (e.g., "spending", "savings", etc.)
	Content          any       `json:"content"`           // Content of the insight, type data jsonb
	ContentText      string    `json:"content_text"`      // Textual representation of the insight content
	ContentEmbedding any       `json:"content_embedding"` // Type data vector for content embedding
	CreatedAt        time.Time `json:"created_at"`        // Timestamp when the insight was created
	UpdatedAt        time.Time `json:"updated_at"`        // Timestamp when the insight was last updated
}
