package models

import "time"

type Receipt struct {
	ReceiptId              string    `json:"receipt_id"`
	UserId                 string    `json:"user_id"`
	FilePath               string    `json:"file_path"`                // Path to the receipt file
	ExtractedReceipt       string    `json:"extracted_receipt"`        // Extracted text from the receipt
	ExtractedTextEmbedding any       `json:"extracted_text_embedding"` // Type data vector for extracted text embedding
	TransactionId          string    `json:"transaction_id"`           // Associated transaction ID
	CreatedAt              time.Time `json:"created_at"`               // Timestamp when the receipt was created
	UpdatedAt              time.Time `json:"updated_at"`               // Timestamp when the receipt was last updated
}
