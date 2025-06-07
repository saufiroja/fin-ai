package log_message

import "github.com/saufiroja/fin-ai/internal/models"

type LogMessageManager interface {
	InsertLogMessage(logMessage *models.LogMessageModel) error
}
