package receipt

import (
	"github.com/saufiroja/fin-ai/internal/contracts/requests"
	"github.com/saufiroja/fin-ai/internal/models"
)

type ReceiptStorer interface {
	InsertReceipt(receipt *models.Receipt) error
	InsertReceiptItem(receiptItem *models.ReceiptItem) error
	GetReceiptsByUserId(userId string) ([]*models.Receipt, error)
	GetDetailReceiptUserById(userId string, receiptId string) (*models.Receipt, error)
	GetReceiptItemsByReceiptId(receiptId string) ([]*models.ReceiptItem, error)
	UpdateReceiptConfirmed(receiptId string, confirmed bool) error
	CountReceiptsByUserId(userId string, req *requests.GetAllReceiptsQuery) (int64, error)
	GetAllReceiptsByUserId(userId string, req *requests.GetAllReceiptsQuery) ([]*models.Receipt, error)
}
