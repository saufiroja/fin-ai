package chat

import "github.com/gofiber/fiber/v2"

type ChatControllerInterface interface {
	CreateChatSession(ctx *fiber.Ctx) error
	FindAllChatSessions(ctx *fiber.Ctx) error
	RenameChatSession(ctx *fiber.Ctx) error
	DeleteChatSession(ctx *fiber.Ctx) error
	SendChatMessage(ctx *fiber.Ctx) error
	GetChatSessionDetail(ctx *fiber.Ctx) error
}
