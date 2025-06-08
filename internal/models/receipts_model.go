package models

import "time"

type Receipt struct {
	ReceiptId                 string    `json:"receipt_id"`
	UserId                    string    `json:"user_id"`
	MerchantName              string    `json:"merchant_name"`
	SubTotal                  int64     `json:"subtotal"`       // total amount before tax
	TotalDiscount             int64     `json:"total_discount"` // total discount applied
	TotalShopping             int64     `json:"total_shopping"` // total amount after discount
	MetaData                  []byte    `json:"meta_data"`      // metadata about the receipt file
	ExtractedReceipt          []byte    `json:"extracted_receipt"`
	ExtractedReceiptEmbedding any       `json:"extracted_receipt_embedding"`
	Confirmed                 bool      `json:"confirmed"`        // whether the receipt is confirmed by the user
	TransactionDate           time.Time `json:"transaction_date"` // date of the transaction
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

type ReceiptItem struct {
	ReceiptItemId  string    `json:"receipt_item_id"`
	ReceiptId      string    `json:"receipt_id"`
	ItemName       string    `json:"item_name"`        // name of the item purchased
	ItemQuantity   int       `json:"item_quantity"`    // quantity of the item purchased
	ItemPrice      int64     `json:"item_price"`       // price of the item
	ItemPriceTotal int64     `json:"item_price_total"` // total price for the item (quantity * item_price)
	ItemDiscount   int64     `json:"item_discount"`    // discount applied to the item
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type MetaData struct {
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"` // size in bytes
	FileType string `json:"file_type"` // e.g., "image/jpeg", "application/pdf"
}
