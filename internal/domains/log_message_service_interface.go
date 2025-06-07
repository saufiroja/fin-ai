package domains

import "github.com/saufiroja/fin-ai/internal/models"

type LogMessageServiceInterface interface {
	InsertLogMessage(logMessage *models.LogMessageModel) error
}
