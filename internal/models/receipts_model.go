package models

import "time"

type Receipt struct {
	ReceiptId                 string    `json:"receipt_id"`
	UserId                    string    `json:"user_id"`
	MerchantName              string    `json:"merchant_name"`
	SubTotal                  int64     `json:"subtotal"`
	TotalDiscount             int64     `json:"total_discount"`
	TotalShopping             int64     `json:"total_shopping"`
	MetaData                  []byte    `json:"-"`
	ExtractedReceipt          []byte    `json:"-"`
	ExtractedReceiptEmbedding any       `json:"-"`
	Confirmed                 bool      `json:"confirmed"`
	TransactionDate           time.Time `json:"transaction_date"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

type ReceiptItem struct {
	ReceiptItemId  string    `json:"receipt_item_id"`
	ReceiptId      string    `json:"receipt_id"`
	ItemName       string    `json:"item_name"`
	ItemQuantity   int       `json:"item_quantity"`
	ItemPrice      int64     `json:"item_price"`
	ItemPriceTotal int64     `json:"item_price_total"`
	ItemDiscount   int64     `json:"item_discount"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type MetaData struct {
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
	FileType string `json:"file_type"`
}
