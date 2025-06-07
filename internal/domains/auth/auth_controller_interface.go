package auth

import "github.com/gofiber/fiber/v2"

type AuthController interface {
	RegisterUser(c *fiber.Ctx) error
	LoginUser(c *fiber.Ctx) error
	LogoutUser(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
}
