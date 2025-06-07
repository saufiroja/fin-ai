package domains

import (
	"context"

	"github.com/saufiroja/fin-ai/internal/models"
)

type ChatServiceInterface interface {
	CreateChatSession(userId string) (*models.ChatSession, error)
	FindAllChatSessions(userId string) ([]*models.ChatSession, error)
	RenameChatSession(chatSession *models.ChatSessionUpdateRequest) error
	DeleteChatSession(chatSessionId, userId string) error
	SendChatMessage(ctx context.Context, message *models.ChatMessageRequest) (*models.ChatMessage, error)
	FindChatSessionDetailByChatSessionIdAndUserId(chatSessionId, userId string) ([]*models.ChatSessionDetail, error)
}
