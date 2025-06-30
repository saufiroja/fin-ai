package receipt

import "github.com/saufiroja/fin-ai/internal/models"

type ReceiptStorer interface {
	InsertReceipt(receipt *models.Receipt) error
	InsertReceiptItem(receiptItem *models.ReceiptItem) error
	GetReceiptsByUserId(userId string) ([]*models.Receipt, error)
	GetDetailReceiptUserById(userId string, receiptId string) (*models.Receipt, error)
	GetReceiptItemsByReceiptId(receiptId string) ([]*models.ReceiptItem, error)
	UpdateReceiptConfirmed(receiptId string, confirmed bool) error
}
