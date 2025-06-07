package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saufiroja/fin-ai/internal/contracts/requests"
	"github.com/saufiroja/fin-ai/internal/contracts/responses"
	"github.com/saufiroja/fin-ai/internal/domains"
)

type userController struct {
	UserService domains.UserServiceInterface
}

func NewUserController(userService domains.UserServiceInterface) UserController {
	return &userController{
		UserService: userService,
	}
}

func (c *userController) UpdateUserById(ctx *fiber.Ctx) error {
	userId := ctx.Params("user_id")
	if userId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.Response{
			Status:  fiber.StatusBadRequest,
			Message: "User ID is required",
		})
	}

	var req requests.UpdateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	err := c.UserService.UpdateUserById(userId, &req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.Response{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.Response{
		Status:  fiber.StatusOK,
		Message: "User updated successfully",
	})
}

func (c *userController) DeleteUserById(ctx *fiber.Ctx) error {
	userId := ctx.Params("user_id")
	if userId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.Response{
			Status:  fiber.StatusBadRequest,
			Message: "User ID is required",
		})
	}

	err := c.UserService.DeleteUserById(userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.Response{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.Response{
		Status:  fiber.StatusOK,
		Message: "user deleted successfully",
	})
}

func (c *userController) GetMe(ctx *fiber.Ctx) error {
	userId := ctx.Locals("user_id").(string)

	if userId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.Response{
			Status:  fiber.StatusBadRequest,
			Message: "User ID not found in token",
		})
	}

	// Ambil user info dari service
	user, err := c.UserService.GetMe(userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to get user information",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.Response{
		Status:  fiber.StatusOK,
		Message: "User information retrieved successfully",
		Data:    user,
	})
}
