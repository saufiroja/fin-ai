package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saufiroja/fin-ai/internal/interfaces"
	"github.com/saufiroja/fin-ai/internal/models"
)

type authController struct {
	authService interfaces.AuthServiceInterface
}

func NewAuthController(authService interfaces.AuthServiceInterface) AuthController {
	return &authController{
		authService: authService,
	}
}

func (c *authController) RegisterUser(ctx *fiber.Ctx) error {
	var req models.User
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	err := c.authService.RegisterUser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to register user",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User registered successfully",
	})
}
