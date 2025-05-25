package constants

import "time"

const (
	HistoryLimit         = 10
	TitleMaxLength       = 50
	DefaultChatTitle     = "New Chat"
	SystemPrompt         = "You are a financial assistant. Provide helpful and accurate responses to user queries."
	TitleSuffix          = "..."
	TitleTruncateLength  = 47
	TitleGenerationModel = "gpt-3.5-turbo" // Lightweight model for title generation
	TitleTimeout         = 30 * time.Second
)
