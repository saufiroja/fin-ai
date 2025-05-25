package interfaces

import "github.com/saufiroja/fin-ai/internal/models"

type LogMessageRepositoryInterface interface {
	InsertLogMessage(logMessage *models.LogMessageModel) error
}
