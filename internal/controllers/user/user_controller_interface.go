package user

import "github.com/gofiber/fiber/v2"

type UserController interface {
	GetMe(c *fiber.Ctx) error
	UpdateUserById(ctx *fiber.Ctx) error
	DeleteUserById(ctx *fiber.Ctx) error
}
