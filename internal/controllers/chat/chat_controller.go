package chat

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saufiroja/fin-ai/internal/interfaces"
	"github.com/saufiroja/fin-ai/internal/models"
)

type chatController struct {
	chatService interfaces.ChatServiceInterface
}

func NewChatController(chatService interfaces.ChatServiceInterface) ChatControllerInterface {
	return &chatController{
		chatService: chatService,
	}
}

func (c *chatController) CreateChatSession(ctx *fiber.Ctx) error {
	userId := ctx.Params("user_id")
	err := c.chatService.CreateChatSession(userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to create chat session",
		})
	}
	return ctx.Status(fiber.StatusCreated).JSON(models.Response{
		Status:  fiber.StatusCreated,
		Message: "Chat session created successfully",
	})
}

func (c *chatController) FindAllChatSessions(ctx *fiber.Ctx) error {
	userId := ctx.Params("user_id")
	chatSessions, err := c.chatService.FindAllChatSessions(userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to find chat sessions",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(models.Response{
		Status:  fiber.StatusOK,
		Message: "Chat sessions retrieved successfully",
		Data:    chatSessions,
	})
}
