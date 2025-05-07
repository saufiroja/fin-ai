package auth

import "github.com/gofiber/fiber/v2"

type AuthController interface {
	RegisterUser(c *fiber.Ctx) error
}
