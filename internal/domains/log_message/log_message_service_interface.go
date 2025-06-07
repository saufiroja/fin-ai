package log_message

import "github.com/saufiroja/fin-ai/internal/models"

type LogMessageService interface {
	InsertLogMessage(logMessage *models.LogMessageModel) error
}
