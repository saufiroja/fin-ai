package interfaces

import "github.com/saufiroja/fin-ai/internal/models"

type ChatRepositoryInterface interface {
	InsertChatSession(chatSession *models.ChatSession) error
	FindAllChatSessions(userId string) ([]*models.ChatSession, error)
}
