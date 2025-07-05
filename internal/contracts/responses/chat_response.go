package responses

import (
	"time"

	"github.com/saufiroja/fin-ai/internal/models"
)

type ChatMessageResponse struct {
	ChatMessageId string          `json:"chat_message_id"`
	ChatSessionId string          `json:"chat_session_id"`
	Conversation  []*Conversation `json:"conversation"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	DeletedAt     time.Time       `json:"deleted_at"`
}

type Conversation struct {
	Sender models.ChatMessageSender `json:"sender"`
	Text   string                   `json:"text"`
}
