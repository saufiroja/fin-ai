package services

import (
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/saufiroja/fin-ai/internal/interfaces"
	"github.com/saufiroja/fin-ai/internal/models"
	logging "github.com/saufiroja/fin-ai/pkg/loggings"
)

type chatService struct {
	chatRepository interfaces.ChatRepositoryInterface
	logging        logging.Logger
}

func NewChatService(chatRepository interfaces.ChatRepositoryInterface, logging logging.Logger) interfaces.ChatServiceInterface {
	return &chatService{
		chatRepository: chatRepository,
		logging:        logging,
	}
}

func (s *chatService) CreateChatSession(userId string) error {
	s.logging.LogInfo(fmt.Sprintf("Creating chat session for user: %s", userId))
	chatSession := &models.ChatSession{
		ChatSessionId: ulid.Make().String(),
		UserId:        userId,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	err := s.chatRepository.InsertChatSession(chatSession)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to create chat session: %s", err.Error()))
		return err
	}

	s.logging.LogInfo("Chat session created successfully")
	return nil
}

func (s *chatService) FindAllChatSessions(userId string) ([]*models.ChatSession, error) {
	s.logging.LogInfo(fmt.Sprintf("Finding all chat sessions for user: %s", userId))
	chatSessions, err := s.chatRepository.FindAllChatSessions(userId)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to find chat sessions: %s", err.Error()))
		return nil, err
	}

	// Return empty slice instead of nil if no chat sessions found
	if chatSessions == nil {
		chatSessions = []*models.ChatSession{}
	}

	s.logging.LogInfo(fmt.Sprintf("Found %d chat sessions", len(chatSessions)))
	return chatSessions, nil
}
