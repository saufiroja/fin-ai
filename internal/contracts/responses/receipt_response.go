package responses

import (
	"time"

	"github.com/saufiroja/fin-ai/internal/models"
)

type ReceiptExtractionResponse struct {
	ExtractedReceipt ExtractedReceiptResponse `json:"extracted_receipt"`
}

type ExtractedReceiptResponse struct {
	MerchantName    string                `json:"merchant_name"`
	SubTotal        int64                 `json:"sub_total"`
	TotalDiscount   int64                 `json:"total_discount"`
	TotalShopping   int64                 `json:"total_shopping"`
	TransactionDate time.Time             `json:"transaction_date"`
	Items           []ReceiptItemResponse `json:"items"`
}

type ReceiptItemResponse struct {
	CategoryId           *string `json:"category_id"`
	ItemName             string  `json:"item_name"`
	ItemQuantity         int     `json:"item_quantity"`
	ItemPrice            int64   `json:"item_price"`
	ItemPriceTotal       int64   `json:"item_price_total"`
	ItemDiscount         int64   `json:"item_discount"`
	AiCategoryConfidence float64 `json:"ai_category_confidence"`
}

type DetailReceiptUserResponse struct {
	ReceiptId       string                `json:"receipt_id"`
	UserId          string                `json:"user_id"`
	MerchantName    string                `json:"merchant_name"`
	SubTotal        int64                 `json:"sub_total"`
	TotalDiscount   int64                 `json:"total_discount"`
	TotalShopping   int64                 `json:"total_shopping"`
	TransactionDate time.Time             `json:"transaction_date"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
	Items           []*models.ReceiptItem `json:"items"`
	Confirmed       bool                  `json:"confirmed"`
}

type ReceiptResponse struct {
	TotalPages  int64             `json:"total_pages"`
	CurrentPage int64             `json:"current_page"`
	Total       int64             `json:"total"`
	Receipts    []*models.Receipt `json:"receipts"`
}
