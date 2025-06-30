package responses

import "time"

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
