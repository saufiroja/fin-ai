package repositories

import (
	"github.com/saufiroja/fin-ai/internal/domains/receipt"
	"github.com/saufiroja/fin-ai/internal/models"
	"github.com/saufiroja/fin-ai/pkg/databases"
)

type receiptRepository struct {
	DB databases.PostgresManager
}

func NewReceiptRepository(db databases.PostgresManager) receipt.ReceiptStorer {
	return &receiptRepository{
		DB: db,
	}
}
func (r *receiptRepository) InsertReceipt(receipt *models.Receipt) error {
	db := r.DB.Connection()

	query := `
    INSERT INTO receipts (
    receipt_id, 
    user_id,
    merchant_name,
    sub_total,
    total_discount,
    total_shopping,
    metadata,
    extracted_receipt, 
    extracted_receipt_embedding, 
    confirmed,
    transaction_date, 
    created_at, 
    updated_at
    ) 
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err := db.Exec(query, receipt.ReceiptId, receipt.UserId, receipt.MerchantName, receipt.SubTotal, receipt.TotalDiscount, receipt.TotalShopping, receipt.MetaData, receipt.ExtractedReceipt, receipt.ExtractedReceiptEmbedding, receipt.Confirmed, receipt.TransactionDate, receipt.CreatedAt, receipt.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *receiptRepository) InsertReceiptItem(receiptItem *models.ReceiptItem) error {
	db := r.DB.Connection()

	query := `
    INSERT INTO receipt_items (
    receipt_item_id, 
    receipt_id, 
    item_name, 
    item_quantity, 
    item_price, 
    item_price_total, 
    item_discount, 
    created_at, 
    updated_at
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())`

	_, err := db.Exec(query, receiptItem.ReceiptItemId, receiptItem.ReceiptId, receiptItem.ItemName, receiptItem.ItemQuantity, receiptItem.ItemPrice, receiptItem.ItemPriceTotal, receiptItem.ItemDiscount)
	if err != nil {
		return err
	}

	return nil
}
