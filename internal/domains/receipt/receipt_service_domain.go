package receipt

import (
	"mime/multipart"

	"github.com/saufiroja/fin-ai/internal/contracts/requests"
	"github.com/saufiroja/fin-ai/internal/contracts/responses"
	"github.com/saufiroja/fin-ai/internal/models"
)

type ReceiptManager interface {
	UploadReceipt(filePath *multipart.FileHeader, userId string) (*models.Receipt, error)
	GetReceiptsByUserId(userId string) ([]*models.Receipt, error)
	GetDetailReceiptUserById(userId string, receiptId string) (*responses.DetailReceiptUserResponse, error)
	UpdateReceiptConfirmed(receiptId string, confirmed bool) error
	GetAllReceiptsByUserId(userId string, req *requests.GetAllReceiptsQuery) (*responses.ReceiptResponse, error)
}
