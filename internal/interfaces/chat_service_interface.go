package interfaces

import "github.com/saufiroja/fin-ai/internal/models"

type ChatServiceInterface interface {
	CreateChatSession(userId string) error
	FindAllChatSessions(userId string) ([]*models.ChatSession, error)
}
