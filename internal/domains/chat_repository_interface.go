package domains

import "github.com/saufiroja/fin-ai/internal/models"

type ChatRepositoryInterface interface {
	InsertChatSession(chatSession *models.ChatSession) error
	FindAllChatSessions(userId string) ([]*models.ChatSession, error)
	RenameChatSession(chatSession *models.ChatSession) error
	DeleteChatSession(chatSessionId, userId string) error
	FindChatSessionByChatSessionIdAndUserId(chatSessionId, userId string) (*models.ChatSession, error)
	InsertChatMessage(chatMessage *models.ChatMessage) error
	FindChatMessagesByChatSessionId(chatSessionId string) ([]*models.ChatMessage, error)
	UpdateChatSessionTitle(chatSession *models.ChatSession) error
	FindChatSessionDetailByChatSessionIdAndUserId(chatSessionId, userId string) ([]*models.ChatSessionDetail, error)
}
