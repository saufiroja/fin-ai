package log_message

import "github.com/saufiroja/fin-ai/internal/models"

type LogMessageRepository interface {
	InsertLogMessage(logMessage *models.LogMessageModel) error
}
