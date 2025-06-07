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
