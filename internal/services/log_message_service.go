package services

import (
	"github.com/saufiroja/fin-ai/internal/domains/log_message"
	"github.com/saufiroja/fin-ai/internal/models"
	logging "github.com/saufiroja/fin-ai/pkg/loggings"
)

type logMessageService struct {
	logMessageRepository log_message.LogMessageStorer
	logging              logging.Logger
}

func NewLogMessageService(
	logMessageRepository log_message.LogMessageStorer,
	logging logging.Logger,
) log_message.LogMessageManager {
	return &logMessageService{
		logMessageRepository: logMessageRepository,
		logging:              logging,
	}
}

func (s *logMessageService) InsertLogMessage(logMessage *models.LogMessage) error {
	s.logging.LogInfo("inserting log message: " + logMessage.Topic)
	err := s.logMessageRepository.InsertLogMessage(logMessage)
	if err != nil {
		s.logging.LogError("failed to insert log message: " + err.Error())
		return err
	}

	s.logging.LogInfo("log message inserted successfully")
	return nil
}
