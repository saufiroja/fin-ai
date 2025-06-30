package receipt

import "github.com/gofiber/fiber/v2"

type ReceiptController interface {
	UploadReceipt(c *fiber.Ctx) error
	GetReceiptsByUserId(c *fiber.Ctx) error
	GetDetailReceiptUserById(c *fiber.Ctx) error
	UpdateReceiptConfirmed(c *fiber.Ctx) error
}
