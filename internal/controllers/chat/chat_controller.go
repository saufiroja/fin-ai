package chat

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saufiroja/fin-ai/internal/interfaces"
	"github.com/saufiroja/fin-ai/internal/models"
	"github.com/saufiroja/fin-ai/internal/utils"
)

type chatController struct {
	chatService interfaces.ChatServiceInterface
	validator   utils.Validator
}

func NewChatController(chatService interfaces.ChatServiceInterface, validator utils.Validator) ChatControllerInterface {
	return &chatController{
		chatService: chatService,
		validator:   validator,
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

func (c *chatController) RenameChatSession(ctx *fiber.Ctx) error {
	chatSession := new(models.ChatSessionUpdateRequest)
	if err := ctx.BodyParser(chatSession); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	if err := c.validator.ValidateStruct(chatSession); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Validation error: " + err.Error(),
		})
	}

	err := c.chatService.RenameChatSession(chatSession)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(models.Response{
		Status:  fiber.StatusOK,
		Message: "Chat session renamed successfully",
	})
}

func (c *chatController) DeleteChatSession(ctx *fiber.Ctx) error {
	chatSessionId := ctx.Params("chat_session_id")
	userId := ctx.Params("user_id")

	err := c.chatService.DeleteChatSession(chatSessionId, userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(models.Response{
		Status:  fiber.StatusOK,
		Message: "Chat session deleted successfully",
	})
}
