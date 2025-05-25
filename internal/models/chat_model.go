package models

import "time"

type ChatSession struct {
	ChatSessionId string    `json:"chat_session_id"`
	UserId        string    `json:"user_id"`
	Title         string    `json:"title"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
	DeletedAt     time.Time `json:"deleted_at,omitempty"`
}

type ChatMessage struct {
	ChatMessageId string            `json:"chat_message_id"`
	ChatSessionId string            `json:"chat_session_id"`
	Message       string            `json:"content"`
	Sender        ChatMessageSender `json:"sender"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	DeletedAt     time.Time         `json:"deleted_at"`
}

type ChatMessageSender string

const (
	ChatMessageSenderUser      ChatMessageSender = "user"
	ChatMessageSenderAssistant ChatMessageSender = "ai"
)

type ChatSessionUpdateRequest struct {
	ChatSessionId string `json:"chat_session_id" validate:"required"`
	UserId        string `json:"user_id" validate:"required"`
	Title         string `json:"title" validate:"required"`
}

type ChatMessageRequest struct {
	ChatSessionId string `json:"chat_session_id" validate:"required"`
	UserId        string `json:"user_id" validate:"required"`
	ModelId       string `json:"model_id" validate:"required"`
	Message       string `json:"message" validate:"required"`
}

type ChatSessionDetail struct {
	ChatMessageId string            `json:"chat_message_id"`
	ChatSessionId string            `json:"chat_session_id"`
	Message       string            `json:"message"`
	Sender        ChatMessageSender `json:"sender"`
	CreatedAt     *time.Time        `json:"created_at"`
	UpdatedAt     *time.Time        `json:"updated_at"`
}
