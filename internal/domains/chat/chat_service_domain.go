package chat

import (
	"context"

	"github.com/saufiroja/fin-ai/internal/contracts/responses"
	"github.com/saufiroja/fin-ai/internal/models"
)

type ChatManager interface {
	CreateChatSession(userId string) (*models.ChatSession, error)
	FindAllChatSessions(userId string) ([]*models.ChatSession, error)
	RenameChatSession(userId, chatSessionId string, chatSession *models.ChatSessionUpdateRequest) error
	DeleteChatSession(chatSessionId, userId string) error
	SendChatMessage(ctx context.Context, message *models.ChatMessageRequest) (*responses.ChatMessageResponse, error)
	FindChatSessionDetailByChatSessionIdAndUserId(chatSessionId, userId string) ([]*models.ChatSessionDetail, error)
}
