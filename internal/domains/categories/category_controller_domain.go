package categories

import "github.com/gofiber/fiber/v2"

type CategoryController interface {
	CreateCategory(c *fiber.Ctx) error
	GetAllCategories(c *fiber.Ctx) error
	UpdateCategoryById(c *fiber.Ctx) error
	DeleteCategoryById(c *fiber.Ctx) error
}
