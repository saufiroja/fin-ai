package responses

import "time"

type ReceiptExtractionResponse struct {
	ExtractedReceipt ExtractedReceiptResponse `json:"extracted_receipt"`
}

type ExtractedReceiptResponse struct {
	MerchantName    string                `json:"merchant_name"`
	SubTotal        int64                 `json:"subtotal"`         // total amount before tax
	TotalDiscount   int64                 `json:"total_discount"`   // total discount applied
	TotalShopping   int64                 `json:"total_shopping"`   // total amount after discount
	TransactionDate time.Time             `json:"transaction_date"` // date of the transaction in ISO 8601 format
	Items           []ReceiptItemResponse `json:"items"`            // list of items purchased
}

type ReceiptItemResponse struct {
	CategoryId           *string `json:"category_id"`
	ItemName             string  `json:"item_name"`        // name of the item purchased
	ItemQuantity         int     `json:"item_quantity"`    // quantity of the item purchased
	ItemPrice            int64   `json:"item_price"`       // price of the item in the smallest currency unit (e.g., cents)
	ItemPriceTotal       int64   `json:"item_price_total"` // total price for the item (quantity * item_price) in the smallest currency unit (e.g., cents)
	ItemDiscount         int64   `json:"item_discount"`    // discount applied to the item in the smallest currency unit (e.g., cents)
	AiCategoryConfidence float64 `json:"ai_category_confidence"`
}
