package auth

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/saufiroja/fin-ai/internal/interfaces"
	"github.com/saufiroja/fin-ai/internal/models"
	"github.com/saufiroja/fin-ai/internal/utils"
)

type authController struct {
	authService interfaces.AuthServiceInterface
	validator   utils.Validator
}

func NewAuthController(authService interfaces.AuthServiceInterface, validator utils.Validator) AuthController {
	return &authController{
		authService: authService,
		validator:   validator,
	}
}

func (c *authController) RegisterUser(ctx *fiber.Ctx) error {
	var req models.RegisterUser
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request",
		})
	}

	if err := c.validator.ValidateStruct(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.Response{
			Status:  fiber.StatusBadRequest,
			Message: fmt.Sprintf("Validation error: %v", err),
		})
	}

	err := c.authService.RegisterUser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to register user",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(models.Response{
		Status:  fiber.StatusOK,
		Message: "user registered successfully",
	})
}

func (c *authController) LoginUser(ctx *fiber.Ctx) error {
	var req models.LoginUser
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request",
		})
	}

	if err := c.validator.ValidateStruct(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.Response{
			Status:  fiber.StatusBadRequest,
			Message: fmt.Sprintf("Validation error: %v", err),
		})
	}

	loginResponse, err := c.authService.LoginUser(&req, ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to login user",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(models.Response{
		Status:  fiber.StatusOK,
		Message: "User logged in successfully",
		Data:    loginResponse,
	})
}

func (c *authController) LogoutUser(ctx *fiber.Ctx) error {
	// Implement logout logic here
	err := c.authService.LogoutUser(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to logout user",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(models.Response{
		Status:  fiber.StatusOK,
		Message: "User logged out successfully",
	})
}

func (c *authController) RefreshToken(ctx *fiber.Ctx) error {
	// get refresh token from cookies
	refreshToken := ctx.Cookies("refresh_token")
	if refreshToken == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(models.Response{
			Status:  fiber.StatusUnauthorized,
			Message: "Refresh token not found",
		})
	}

	// validate refresh token
	jwtToken, err := c.authService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(models.Response{
			Status:  fiber.StatusUnauthorized,
			Message: "Invalid refresh token",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(models.Response{
		Status:  fiber.StatusOK,
		Message: "Access token refreshed successfully",
		Data:    jwtToken,
	})
}
