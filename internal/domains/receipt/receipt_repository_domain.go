package receipt

import "github.com/saufiroja/fin-ai/internal/models"

type ReceiptStorer interface {
	InsertReceipt(receipt *models.Receipt) error
	InsertReceiptItem(receiptItem *models.ReceiptItem) error
}
