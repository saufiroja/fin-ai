package receipt

import "mime/multipart"

type ReceiptManager interface {
	UploadReceipt(filePath *multipart.FileHeader, userId string) error
}
