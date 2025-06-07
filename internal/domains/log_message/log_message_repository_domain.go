package log_message

import "github.com/saufiroja/fin-ai/internal/models"

type LogMessageStorer interface {
	InsertLogMessage(logMessage *models.LogMessageModel) error
}
