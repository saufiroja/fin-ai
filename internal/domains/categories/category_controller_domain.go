package categories

import "github.com/gofiber/fiber/v2"

type CategoryController interface {
	CreateCategory(c *fiber.Ctx) error
}
