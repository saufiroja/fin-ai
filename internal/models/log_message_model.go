package models

import "time"

type LogMessage struct {
	LogMessageId string    `json:"log_message_id"`
	UserId       string    `json:"user_id"`
	Message      string    `json:"message"`
	Response     string    `json:"response"`
	InputToken   int       `json:"input_token"`
	OutputToken  int       `json:"output_token"`
	Topic        string    `json:"topic"`
	Model        string    `json:"model"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
