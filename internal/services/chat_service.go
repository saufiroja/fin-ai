package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/openai/openai-go"
	"github.com/saufiroja/fin-ai/internal/interfaces"
	"github.com/saufiroja/fin-ai/internal/models"
	"github.com/saufiroja/fin-ai/pkg/llm"
	logging "github.com/saufiroja/fin-ai/pkg/loggings"
)

type chatService struct {
	chatRepository interfaces.ChatRepositoryInterface
	logging        logging.Logger
	llmClient      llm.OpenAI
}

func NewChatService(chatRepository interfaces.ChatRepositoryInterface, logging logging.Logger, llmClient llm.OpenAI) interfaces.ChatServiceInterface {
	return &chatService{
		chatRepository: chatRepository,
		logging:        logging,
		llmClient:      llmClient,
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

func (s *chatService) RenameChatSession(req *models.ChatSessionUpdateRequest) error {
	s.logging.LogInfo(fmt.Sprintf("Renaming chat session: %s", req.ChatSessionId))

	_, err := s.chatRepository.FindChatSessionByChatSessionIdAndUserId(req.ChatSessionId, req.UserId)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Chat session not found: %s", err.Error()))
		return errors.New("chat session not found")
	}

	chatSession := &models.ChatSession{
		ChatSessionId: req.ChatSessionId,
		UserId:        req.UserId,
		Title:         req.Title,
		UpdatedAt:     time.Now(),
	}

	err = s.chatRepository.RenameChatSession(chatSession)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to rename chat session: %s", err.Error()))
		return errors.New("failed to rename chat session")
	}

	s.logging.LogInfo("Chat session renamed successfully")
	return nil
}

func (s *chatService) DeleteChatSession(chatSessionId, userId string) error {
	s.logging.LogInfo(fmt.Sprintf("Deleting chat session: %s", chatSessionId))

	_, err := s.chatRepository.FindChatSessionByChatSessionIdAndUserId(chatSessionId, userId)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Chat session not found: %s", err.Error()))
		return errors.New("chat session not found")
	}

	err = s.chatRepository.DeleteChatSession(chatSessionId, userId)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to delete chat session: %s", err.Error()))
		return errors.New("failed to delete chat session")
	}

	s.logging.LogInfo("Chat session deleted successfully")
	return nil
}

func (s *chatService) SendChatMessage(ctx context.Context, message *models.ChatMessageRequest) (*models.ChatMessage, error) {
	s.logging.LogInfo(fmt.Sprintf("Sending chat message to session: %s", message.ChatSessionId))

	// Validate chat session exists
	chatSession, err := s.chatRepository.FindChatSessionByChatSessionIdAndUserId(message.ChatSessionId, message.UserId)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Chat session not found: %s", err.Error()))
		return nil, errors.New("chat session not found")
	}

	if chatSession == nil {
		s.logging.LogError("Chat session is nil")
		return nil, errors.New("chat session not found")
	}

	// Create and save user message
	userMessage := &models.ChatMessage{
		ChatMessageId: ulid.Make().String(),
		ChatSessionId: chatSession.ChatSessionId,
		Message:       message.Message,
		Sender:        models.ChatMessageSenderUser,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err = s.chatRepository.InsertChatMessage(userMessage)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to save user message: %s", err.Error()))
		return nil, errors.New("failed to save user message")
	}

	// Get conversation history for context
	conversationHistory, err := s.chatRepository.FindChatMessagesByChatSessionId(message.ChatSessionId)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to get conversation history: %s", err.Error()))
		// Continue without history if this fails
		conversationHistory = []*models.ChatMessage{}
	}

	// Build messages for LLM
	messages := s.buildLLMMessages(conversationHistory, message.Message)

	// Get AI response
	aiResponse, err := s.llmClient.SendChat(ctx, messages)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to get AI response: %s", err.Error()))
		return userMessage, errors.New("failed to get AI response")
	}

	// Create and save AI response message
	aiMessage := &models.ChatMessage{
		ChatMessageId: ulid.Make().String(),
		ChatSessionId: chatSession.ChatSessionId,
		Message:       aiResponse,
		Sender:        models.ChatMessageSenderAssistant,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err = s.chatRepository.InsertChatMessage(aiMessage)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to save AI message: %s", err.Error()))
		// Return user message even if AI message save fails
		return userMessage, nil
	}

	// Update chat session title if it's the first message
	if len(conversationHistory) == 0 && chatSession.Title == "" {
		go s.updateChatSessionTitle(chatSession.ChatSessionId, message.Message)
	}

	s.logging.LogInfo("Chat message sent successfully")
	return aiMessage, nil
}

func (s *chatService) buildLLMMessages(history []*models.ChatMessage, currentMessage string) []openai.ChatCompletionMessageParamUnion {
	var messages []openai.ChatCompletionMessageParamUnion

	// Add system message for context
	messages = append(messages,
		openai.SystemMessage("You are a financial assistant. Provide helpful and accurate responses to user queries."))

	// Add conversation history (limit to last 10 messages for context)
	historyLimit := 10
	start := 0
	if len(history) > historyLimit {
		start = len(history) - historyLimit
	}

	for i := start; i < len(history); i++ {
		msg := history[i]

		if msg.Sender == models.ChatMessageSenderAssistant {
			messages = append(messages, openai.AssistantMessage(msg.Message))
		} else {
			userMsg := openai.UserMessage(msg.Message)
			messages = append(messages, userMsg)
		}
	}

	// Add current user message
	currentUserMsg := openai.UserMessage(currentMessage)

	messages = append(messages, currentUserMsg)
	fmt.Printf("Built %d messages for LLM\n", len(messages))

	return messages
}

// updateChatSessionTitle updates the chat session title based on the first message
func (s *chatService) updateChatSessionTitle(chatSessionId, firstMessage string) {
	// Generate a short title from the first message (first 50 characters)
	title := firstMessage
	if len(title) > 50 {
		title = title[:47] + "..."
	}

	chatSession := &models.ChatSession{
		ChatSessionId: chatSessionId,
		Title:         title,
		UpdatedAt:     time.Now(),
	}

	err := s.chatRepository.UpdateChatSessionTitle(chatSession)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to update chat session title: %s", err.Error()))
	}
}

func (s *chatService) FindChatSessionDetailByChatSessionIdAndUserId(chatSessionId, userId string) ([]*models.ChatSessionDetail, error) {
	chatSessionDetail, err := s.chatRepository.FindChatSessionDetailByChatSessionIdAndUserId(chatSessionId, userId)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to find chat session detail: %s", err.Error()))
		return nil, errors.New("failed to find chat session detail")
	}

	return chatSessionDetail, nil
}
