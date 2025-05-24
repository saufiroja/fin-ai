package chat

import "github.com/gofiber/fiber/v2"

type ChatControllerInterface interface {
	CreateChatSession(ctx *fiber.Ctx) error
	FindAllChatSessions(ctx *fiber.Ctx) error
}
