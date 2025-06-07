package models

import "time"

type OcrCache struct {
	OcrCacheId      string    `json:"ocr_cache_id"`
	FileHash        string    `json:"file_hash"`
	ExtractedData   any       `json:"extracted_data"`   // type data jsonb for extracted data
	ConfidanceScore float64   `json:"confidence_score"` // Confidence score of the OCR extraction
	CreatedAt       time.Time `json:"created_at"`       // Timestamp of when the OCR cache was created
	UpdatedAt       time.Time `json:"updated_at"`       // Timestamp of when the OCR cache was last updated
}
