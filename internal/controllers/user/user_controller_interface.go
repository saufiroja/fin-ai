package user

import "github.com/gofiber/fiber/v2"

type UserControllerInterface interface {
	GetMe(c *fiber.Ctx) error
	UpdateUserById(ctx *fiber.Ctx) error
	DeleteUserById(ctx *fiber.Ctx) error
}
