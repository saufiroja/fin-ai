package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saufiroja/fin-ai/internal/interfaces"
	"github.com/saufiroja/fin-ai/internal/models"
)

type userController struct {
	UserService interfaces.UserServiceInterface
}

func NewUserController(userService interfaces.UserServiceInterface) UserControllerInterface {
	return &userController{
		UserService: userService,
	}
}

func (c *userController) UpdateUserById(ctx *fiber.Ctx) error {
	userId := ctx.Params("user_id")
	if userId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.Response{
			Status:  fiber.StatusBadRequest,
			Message: "User ID is required",
		})
	}

	var req models.UpdateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	err := c.UserService.UpdateUserById(userId, &req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(models.Response{
		Status:  fiber.StatusOK,
		Message: "User updated successfully",
	})
}

func (c *userController) DeleteUserById(ctx *fiber.Ctx) error {
	userId := ctx.Params("user_id")
	if userId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.Response{
			Status:  fiber.StatusBadRequest,
			Message: "User ID is required",
		})
	}

	err := c.UserService.DeleteUserById(userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(models.Response{
		Status:  fiber.StatusOK,
		Message: "user deleted successfully",
	})
}

func (c *userController) GetMe(ctx *fiber.Ctx) error {
	userId := ctx.Locals("user_id").(string)

	if userId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.Response{
			Status:  fiber.StatusBadRequest,
			Message: "User ID not found in token",
		})
	}

	// Ambil user info dari service
	user, err := c.UserService.GetMe(userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to get user information",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(models.Response{
		Status:  fiber.StatusOK,
		Message: "User information retrieved successfully",
		Data:    user,
	})
}
