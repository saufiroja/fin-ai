package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/saufiroja/fin-ai/internal/contracts/responses"
	"github.com/saufiroja/fin-ai/internal/domains/chat"
	"github.com/saufiroja/fin-ai/internal/models"
	"github.com/saufiroja/fin-ai/internal/utils"
)

type chatController struct {
	chatService chat.ChatManager
	validator   utils.Validator
}

func NewChatController(chatService chat.ChatManager, validator utils.Validator) chat.ChatController {
	return &chatController{
		chatService: chatService,
		validator:   validator,
	}
}

func (c *chatController) CreateChatSession(ctx *fiber.Ctx) error {
	userId := ctx.Locals("user_id").(string)
	chatSession, err := c.chatService.CreateChatSession(userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to create chat session",
		})
	}
	return ctx.Status(fiber.StatusCreated).JSON(responses.Response{
		Status:  fiber.StatusCreated,
		Message: "Chat session created successfully",
		Data:    chatSession,
	})
}

func (c *chatController) FindAllChatSessions(ctx *fiber.Ctx) error {
	userId := ctx.Locals("user_id").(string)
	chatSessions, err := c.chatService.FindAllChatSessions(userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to find chat sessions",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(responses.Response{
		Status:  fiber.StatusOK,
		Message: "Chat sessions retrieved successfully",
		Data:    chatSessions,
	})
}

func (c *chatController) RenameChatSession(ctx *fiber.Ctx) error {
	chatSession := new(models.ChatSessionUpdateRequest)
	if err := ctx.BodyParser(chatSession); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	if err := c.validator.ValidateStruct(chatSession); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Validation error: " + err.Error(),
		})
	}

	err := c.chatService.RenameChatSession(chatSession)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.Response{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(responses.Response{
		Status:  fiber.StatusOK,
		Message: "Chat session renamed successfully",
	})
}

func (c *chatController) DeleteChatSession(ctx *fiber.Ctx) error {
	chatSessionId := ctx.Params("chat_session_id")
	userId := ctx.Locals("user_id").(string)

	err := c.chatService.DeleteChatSession(chatSessionId, userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.Response{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.Response{
		Status:  fiber.StatusOK,
		Message: "Chat session deleted successfully",
	})
}

func (c *chatController) SendChatMessage(ctx *fiber.Ctx) error {
	message := new(models.ChatMessageRequest)
	message.UserId = ctx.Locals("user_id").(string)

	if err := ctx.BodyParser(message); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}
	fmt.Println("Received message:", message)

	if err := c.validator.ValidateStruct(message); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Validation error: " + err.Error(),
		})
	}

	response, err := c.chatService.SendChatMessage(ctx.Context(), message)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to send chat message",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.Response{
		Status:  fiber.StatusOK,
		Message: "Chat message send successfully",
		Data:    response,
	})
}

func (c *chatController) GetChatSessionDetail(ctx *fiber.Ctx) error {
	chatSessionId := ctx.Params("chat_session_id")
	userId := ctx.Locals("user_id").(string)

	messages, err := c.chatService.FindChatSessionDetailByChatSessionIdAndUserId(chatSessionId, userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to retrieve chat session details",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.Response{
		Status:  fiber.StatusOK,
		Message: "Chat session details retrieved successfully",
		Data:    messages,
	})
}
