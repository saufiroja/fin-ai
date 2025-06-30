package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/saufiroja/fin-ai/internal/constants/prompt"
	"github.com/saufiroja/fin-ai/internal/domains/categories"
	"github.com/saufiroja/fin-ai/internal/domains/chat"
	"github.com/saufiroja/fin-ai/internal/domains/log_message"
	"github.com/saufiroja/fin-ai/internal/domains/model_registry"
	"github.com/saufiroja/fin-ai/internal/domains/transaction"
	"github.com/saufiroja/fin-ai/internal/models"
	"github.com/saufiroja/fin-ai/pkg/llm"
	logging "github.com/saufiroja/fin-ai/pkg/loggings"
	"google.golang.org/genai"
)

type chatService struct {
	chatRepository     chat.ChatStorer
	logging            logging.Logger
	geminiClient       llm.Gemini
	modelRegistry      model_registry.ModelRegistryStorer
	logMessageService  log_message.LogMessageManager
	transactionService transaction.TransactionManager
	categoryService    categories.CategoryManager
}

func NewChatService(
	chatRepository chat.ChatStorer,
	logging logging.Logger,
	geminiClient llm.Gemini,
	modelRegistry model_registry.ModelRegistryStorer,
	logMessageService log_message.LogMessageManager,
	transactionService transaction.TransactionManager,
	categoryService categories.CategoryManager,
) chat.ChatManager {
	// Set transaction service to gemini client
	geminiClient.SetTransactionService(transactionService)
	// Set category service to gemini client
	geminiClient.SetCategoryService(categoryService)

	return &chatService{
		chatRepository:     chatRepository,
		logging:            logging,
		geminiClient:       geminiClient,
		modelRegistry:      modelRegistry,
		logMessageService:  logMessageService,
		transactionService: transactionService,
		categoryService:    categoryService,
	}
}

func (s *chatService) CreateChatSession(userId string) (*models.ChatSession, error) {
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
		return nil, err
	}

	s.logging.LogInfo("Chat session created successfully")
	return chatSession, nil
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

func (s *chatService) FindChatSessionDetailByChatSessionIdAndUserId(chatSessionId, userId string) ([]*models.ChatSessionDetail, error) {
	chatSessionDetail, err := s.chatRepository.FindChatSessionDetailByChatSessionIdAndUserId(chatSessionId, userId)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to find chat session detail: %s", err.Error()))
		return nil, errors.New("failed to find chat session detail")
	}

	return chatSessionDetail, nil
}

// GetSupportedModes returns a list of supported chat modes
func (s *chatService) GetSupportedModes() []models.Mode {
	return []models.Mode{
		models.ModeChat,
		models.ModeAgent,
	}
}

// GetModeDescription returns description for each mode
func (s *chatService) GetModeDescription(mode models.Mode) string {
	switch mode {
	case models.ModeChat:
		return "Ask mode: AI responds to questions and provides helpful information"
	case models.ModeAgent:
		return "Agent mode: AI proactively analyzes data and provides insights and recommendations"
	default:
		return "Unknown mode"
	}
}

// validateMode validates if the provided mode is supported
func (s *chatService) validateMode(mode models.Mode) error {
	switch mode {
	case models.ModeChat, models.ModeAgent:
		return nil
	default:
		return fmt.Errorf("unsupported mode: %s", mode)
	}
}

// getSystemPromptByMode returns the appropriate system prompt based on the mode
func (s *chatService) getSystemPromptByMode(mode models.Mode) string {
	switch mode {
	case models.ModeAgent:
		return prompt.ChatAgentSystemPrompt
	case models.ModeChat:
		fallthrough
	default:
		return prompt.ChatSystemPrompt
	}
}

func (s *chatService) SendChatMessage(ctx context.Context, req *models.ChatMessageRequest) (*models.ChatMessage, error) {
	// Set default mode if empty
	if req.Mode == "" {
		req.Mode = models.ModeChat
	}

	// Validate mode
	if err := s.validateMode(req.Mode); err != nil {
		s.logging.LogError(fmt.Sprintf("Invalid mode provided: %s", err.Error()))
		return nil, err
	}

	s.logging.LogInfo(fmt.Sprintf("Processing chat message in %s mode for session: %s", req.Mode, req.ChatSessionId))

	var responseAi *models.ChatMessage

	// Handle different modes
	switch req.Mode {
	case models.ModeAgent:
		// Use RunAgent for agent mode
		s.logging.LogInfo("Using RunAgent for agent mode")
		response, err := s.geminiClient.RunAgent(ctx, req.Message, req.UserId)
		if err != nil {
			s.logging.LogError(fmt.Sprintf("Failed to run Gemini agent: %s", err.Error()))
			return nil, fmt.Errorf("failed to run Gemini agent: %w", err)
		}

		responseAi = &models.ChatMessage{
			ChatSessionId: req.ChatSessionId,
			ChatMessageId: "",
			Message:       response.Response.(string),
		}

	case models.ModeChat:
		fallthrough
	default:
		// Use regular Run for chat mode
		s.logging.LogInfo("Using Run for chat mode")

		// Get appropriate system prompt based on mode
		systemPrompt := s.getSystemPromptByMode(req.Mode)

		partsAi := []*genai.Part{
			genai.NewPartFromText(systemPrompt),
		}
		parts := []*genai.Part{
			genai.NewPartFromText(req.Message),
		}

		message := []*genai.Content{
			genai.NewContentFromParts(partsAi, genai.RoleModel),
			genai.NewContentFromParts(parts, genai.RoleUser),
		}

		response, err := s.geminiClient.Run(ctx, "gemini-2.5-flash", message)
		if err != nil {
			s.logging.LogError(fmt.Sprintf("Failed to run Gemini client: %s", err.Error()))
			return nil, fmt.Errorf("failed to run Gemini client: %w", err)
		}

		responseAi = &models.ChatMessage{
			ChatSessionId: req.ChatSessionId,
			ChatMessageId: "",
			Message:       response.Response.(string),
		}
	}

	s.logging.LogInfo("Chat message processed successfully")
	return responseAi, nil
}
