package receipt

import "github.com/gofiber/fiber/v2"

type ReceiptController interface {
	UploadReceipt(c *fiber.Ctx) error
}
