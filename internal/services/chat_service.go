package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/openai/openai-go"
	"github.com/saufiroja/fin-ai/internal/constants"
	"github.com/saufiroja/fin-ai/internal/domains/chat"
	"github.com/saufiroja/fin-ai/internal/domains/log_message"
	"github.com/saufiroja/fin-ai/internal/domains/model_registry"
	"github.com/saufiroja/fin-ai/internal/models"
	"github.com/saufiroja/fin-ai/pkg/llm"
	logging "github.com/saufiroja/fin-ai/pkg/loggings"
)

type chatService struct {
	chatRepository    chat.ChatRepository
	logging           logging.Logger
	llmClient         llm.OpenAI
	modelRegistry     model_registry.ModelRegistryRepository
	logMessageService log_message.LogMessageService
}

func NewChatService(
	chatRepository chat.ChatRepository,
	logging logging.Logger,
	llmClient llm.OpenAI,
	modelRegistry model_registry.ModelRegistryRepository,
	logMessageService log_message.LogMessageService,
) chat.ChatService {
	return &chatService{
		chatRepository:    chatRepository,
		logging:           logging,
		llmClient:         llmClient,
		modelRegistry:     modelRegistry,
		logMessageService: logMessageService,
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

func (s *chatService) SendChatMessage(ctx context.Context, req *models.ChatMessageRequest) (*models.ChatMessage, error) {
	s.logging.LogInfo(fmt.Sprintf("Processing chat message for session: %s", req.ChatSessionId))

	// Validate and get model
	model, err := s.validateAndGetModel(req.ModelId)
	if err != nil {
		return nil, err
	}

	// Validate and get chat session
	chatSession, err := s.validateAndGetChatSession(req.ChatSessionId, req.UserId)
	if err != nil {
		return nil, err
	}

	// Get conversation history (non-blocking failure)
	conversationHistory := s.getConversationHistory(req.ChatSessionId)

	// Save user message
	userMessage, err := s.createAndSaveUserMessage(req, chatSession.ChatSessionId)
	if err != nil {
		return nil, err
	}

	// Get AI response
	aiMessage, err := s.processAIResponse(ctx, model, conversationHistory, req.Message, chatSession.ChatSessionId)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to get AI response: %s", err.Error()))
		return userMessage, fmt.Errorf("failed to get AI response: %w", err)
	}

	// Handle async operations
	s.handleAsyncOperations(chatSession, conversationHistory, req, aiMessage, model)

	s.logging.LogInfo("Chat message processed successfully")
	return aiMessage, nil
}

func (s *chatService) validateAndGetModel(modelId string) (*models.ModelRegistry, error) {
	model, err := s.modelRegistry.FindModelById(modelId)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to find model: %s", err.Error()))
		return nil, fmt.Errorf("failed to find model: %w", err)
	}

	if model == nil {
		s.logging.LogError("Model not found in registry")
		return nil, fmt.Errorf("model not found in registry")
	}

	return model, nil
}

func (s *chatService) validateAndGetChatSession(chatSessionId, userId string) (*models.ChatSession, error) {
	chatSession, err := s.chatRepository.FindChatSessionByChatSessionIdAndUserId(chatSessionId, userId)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Chat session not found: %s", err.Error()))
		return nil, fmt.Errorf("chat session not found: %w", err)
	}

	if chatSession == nil {
		s.logging.LogError("Chat session is nil")
		return nil, fmt.Errorf("chat session not found")
	}

	return chatSession, nil
}

func (s *chatService) getConversationHistory(chatSessionId string) []*models.ChatMessage {
	conversationHistory, err := s.chatRepository.FindChatMessagesByChatSessionId(chatSessionId)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to get conversation history: %s", err.Error()))
		return []*models.ChatMessage{}
	}
	return conversationHistory
}

func (s *chatService) createAndSaveUserMessage(req *models.ChatMessageRequest, chatSessionId string) (*models.ChatMessage, error) {
	userMessage := &models.ChatMessage{
		ChatMessageId: ulid.Make().String(),
		ChatSessionId: chatSessionId,
		Message:       req.Message,
		Sender:        models.ChatMessageSenderUser,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.chatRepository.InsertChatMessage(userMessage); err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to save user message: %s", err.Error()))
		return nil, fmt.Errorf("failed to save user message: %w", err)
	}

	return userMessage, nil
}

func (s *chatService) processAIResponse(ctx context.Context, model *models.ModelRegistry, history []*models.ChatMessage, currentMessage, chatSessionId string) (*models.ChatMessage, error) {
	// Build messages for LLM
	messages := s.buildLLMMessages(history, currentMessage)

	// Get AI response
	aiResponse, err := s.llmClient.SendChat(ctx, model.Name, messages)
	if err != nil {
		return nil, fmt.Errorf("LLM request failed: %w", err)
	}

	// Create and save AI response message
	aiMessage := &models.ChatMessage{
		ChatMessageId: ulid.Make().String(),
		ChatSessionId: chatSessionId,
		Message:       aiResponse.Response,
		Sender:        models.ChatMessageSenderAssistant,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.chatRepository.InsertChatMessage(aiMessage); err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to save AI message: %s", err.Error()))
		// Continue execution even if save fails
	}

	return aiMessage, nil
}

func (s *chatService) handleAsyncOperations(chatSession *models.ChatSession, conversationHistory []*models.ChatMessage, req *models.ChatMessageRequest, aiMessage *models.ChatMessage, model *models.ModelRegistry) {
	var wg sync.WaitGroup

	// Update chat session title if needed (async)
	if s.shouldUpdateTitle(conversationHistory, chatSession.Title) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.updateChatSessionTitle(chatSession.ChatSessionId, chatSession.UserId, req.Message)
		}()
	}

	modelName := model.Name

	// Log chat message (async)
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.logChatMessage(req, aiMessage, modelName)
	}()

	// Optional: Wait for all async operations to complete in a separate goroutine
	go func() {
		wg.Wait()
		s.logging.LogInfo("All async operations completed")
	}()
}

func (s *chatService) shouldUpdateTitle(conversationHistory []*models.ChatMessage, currentTitle string) bool {
	return len(conversationHistory) == 0 && currentTitle == constants.DefaultChatTitle
}

func (s *chatService) buildLLMMessages(history []*models.ChatMessage, currentMessage string) []openai.ChatCompletionMessageParamUnion {
	// Calculate total capacity needed
	totalMessages := 1 + // system message
		s.min(len(history), constants.HistoryLimit) + // limited history
		1 // current message

	messages := make([]openai.ChatCompletionMessageParamUnion, 0, totalMessages)

	// Add system message
	messages = append(messages, openai.SystemMessage(constants.SystemPrompt))

	// Add limited conversation history
	start := s.max(0, len(history)-constants.HistoryLimit)
	for i := start; i < len(history); i++ {
		msg := history[i]
		switch msg.Sender {
		case models.ChatMessageSenderAssistant:
			messages = append(messages, openai.AssistantMessage(msg.Message))
		case models.ChatMessageSenderUser:
			messages = append(messages, openai.UserMessage(msg.Message))
		}
	}

	// Add current user message
	messages = append(messages, openai.UserMessage(currentMessage))

	s.logging.LogInfo(fmt.Sprintf("Built %d messages for LLM", len(messages)))
	return messages
}

func (s *chatService) updateChatSessionTitle(chatSessionId, userId, firstMessage string) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.TitleTimeout)
	defer cancel()

	title := s.generateTitleWithLLM(ctx, firstMessage)

	// Fallback to simple truncation if LLM fails
	if title == "" {
		title = s.generateSimpleTitle(firstMessage)
	}

	chatSession := &models.ChatSession{
		ChatSessionId: chatSessionId,
		UserId:        userId,
		Title:         title,
		UpdatedAt:     time.Now(),
	}

	if err := s.chatRepository.UpdateChatSessionTitle(chatSession); err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to update chat session title: %s", err.Error()))
	} else {
		s.logging.LogInfo(fmt.Sprintf("Updated chat session title: %s", title))
	}
}

func (s *chatService) generateTitleWithLLM(ctx context.Context, message string) string {
	s.logging.LogInfo("Generating title using LLM")

	// Create prompt for title generation
	titlePrompt := fmt.Sprintf(`Generate a concise, descriptive title (max 50 characters) for a conversation that starts with this message:

"%s"

Requirements:
- Maximum 50 characters
- Clear and descriptive
- No quotes or special formatting
- Summarize the main topic or intent
- Professional tone

Title:`, message)

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage("You are a helpful assistant that creates concise, descriptive titles for conversations. Respond with only the title, no additional text."),
		openai.UserMessage(titlePrompt),
	}

	// Use a lightweight model for title generation
	response, err := s.llmClient.SendChat(ctx, constants.TitleGenerationModel, messages)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to generate title with LLM: %s", err.Error()))
		return ""
	}

	// Clean and validate the response
	title := s.cleanGeneratedTitle(response.Response)

	s.logging.LogInfo(fmt.Sprintf("Generated title via LLM: %s", title))
	return title
}

func (s *chatService) cleanGeneratedTitle(rawTitle string) string {
	// Remove quotes if present
	title := strings.Trim(rawTitle, `"'`)

	// Remove any newlines or excessive whitespace
	title = strings.TrimSpace(title)
	title = strings.ReplaceAll(title, "\n", " ")

	// Ensure it doesn't exceed max length
	if len(title) > constants.TitleMaxLength {
		title = title[:constants.TitleTruncateLength] + constants.TitleSuffix
	}

	// Ensure it's not empty
	if title == "" {
		return constants.DefaultChatTitle
	}

	return title
}

func (s *chatService) generateSimpleTitle(message string) string {
	if len(message) <= constants.TitleMaxLength {
		return message
	}
	return message[:constants.TitleTruncateLength] + constants.TitleSuffix
}

func (s *chatService) logChatMessage(req *models.ChatMessageRequest, aiMessage *models.ChatMessage, modelName string) {
	s.logging.LogInfo("Logging chat message")

	logMessage := &models.LogMessageModel{
		LogMessageId: ulid.Make().String(),
		UserId:       req.UserId,
		Message:      req.Message,
		Response:     aiMessage.Message,
		InputToken:   0, // You might want to get this from aiResponse
		OutputToken:  0, // You might want to get this from aiResponse
		Topic:        "chat",
		Model:        modelName,
	}

	if err := s.logMessageService.InsertLogMessage(logMessage); err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to log chat message: %s", err.Error()))
	} else {
		s.logging.LogInfo("Chat message logged successfully")
	}
}

// Helper functions for Go versions that don't have min/max built-ins
func (s *chatService) min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (s *chatService) max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
