package constants

import (
	"time"

	"github.com/saufiroja/fin-ai/internal/constants/prompt"
)

const (
	HistoryLimit         = 10
	TitleMaxLength       = 50
	DefaultChatTitle     = "New Chat"
	TitleSuffix          = "..."
	TitleTruncateLength  = 47
	TitleGenerationModel = "gpt-3.5-turbo" // Lightweight model for title generation
	TitleTimeout         = 30 * time.Second
)

// SystemPrompt is the main system prompt for chat - moved to prompt package
// Keeping this for backward compatibility
var SystemPrompt = prompt.ChatSystemPrompt
