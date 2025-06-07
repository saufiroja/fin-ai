package services

import (
	"github.com/saufiroja/fin-ai/internal/domains"
	"github.com/saufiroja/fin-ai/internal/models"
	logging "github.com/saufiroja/fin-ai/pkg/loggings"
)

type logMessageService struct {
	logMessageRepository domains.LogMessageRepositoryInterface
	logging              logging.Logger
}

func NewLogMessageService(logMessageRepository domains.LogMessageRepositoryInterface, logging logging.Logger) domains.LogMessageServiceInterface {
	return &logMessageService{
		logMessageRepository: logMessageRepository,
		logging:              logging,
	}
}

func (s *logMessageService) InsertLogMessage(logMessage *models.LogMessageModel) error {
	s.logging.LogInfo("inserting log message: " + logMessage.Message)
	err := s.logMessageRepository.InsertLogMessage(logMessage)
	if err != nil {
		s.logging.LogError("failed to insert log message: " + err.Error())
		return err
	}

	s.logging.LogInfo("log message inserted successfully")
	return nil
}
