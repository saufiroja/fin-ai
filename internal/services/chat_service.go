package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	"github.com/saufiroja/fin-ai/internal/constants/prompt"
	"github.com/saufiroja/fin-ai/internal/contracts/requests"
	"github.com/saufiroja/fin-ai/internal/contracts/responses"
	"github.com/saufiroja/fin-ai/internal/domains/categories"
	"github.com/saufiroja/fin-ai/internal/domains/chat"
	"github.com/saufiroja/fin-ai/internal/domains/log_message"
	"github.com/saufiroja/fin-ai/internal/domains/model_registry"
	"github.com/saufiroja/fin-ai/internal/domains/receipt"
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
	openaiClient       llm.OpenAI
	modelRegistry      model_registry.ModelRegistryStorer
	logMessageService  log_message.LogMessageManager
	transactionService transaction.TransactionManager
	categoryService    categories.CategoryManager
	receiptService     receipt.ReceiptManager
}

func NewChatService(
	chatRepository chat.ChatStorer,
	logging logging.Logger,
	geminiClient llm.Gemini,
	openaiClient llm.OpenAI,
	modelRegistry model_registry.ModelRegistryStorer,
	logMessageService log_message.LogMessageManager,
	transactionService transaction.TransactionManager,
	categoryService categories.CategoryManager,
	receiptService receipt.ReceiptManager,
) chat.ChatManager {
	// Set transaction service to gemini client
	geminiClient.SetTransactionService(transactionService)
	// Set category service to gemini client
	geminiClient.SetCategoryService(categoryService)

	return &chatService{
		chatRepository:     chatRepository,
		logging:            logging,
		geminiClient:       geminiClient,
		openaiClient:       openaiClient,
		modelRegistry:      modelRegistry,
		logMessageService:  logMessageService,
		transactionService: transactionService,
		categoryService:    categoryService,
		receiptService:     receiptService,
	}
}

func (s *chatService) CreateChatSession(userId string) (*models.ChatSession, error) {
	s.logging.LogInfo(fmt.Sprintf("Creating chat session for user: %s", userId))
	chatSession := &models.ChatSession{
		ChatSessionId: ulid.Make().String(),
		UserId:        userId,
		Title:         "New Chat",
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

func (s *chatService) RenameChatSession(userId, chatSessionId string, req *models.ChatSessionUpdateRequest) error {
	s.logging.LogInfo(fmt.Sprintf("Renaming chat session: %s", chatSessionId))

	_, err := s.chatRepository.FindChatSessionByChatSessionIdAndUserId(chatSessionId, userId)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Chat session not found: %s", err.Error()))
		return errors.New("chat session not found")
	}

	chatSession := &models.ChatSession{
		ChatSessionId: chatSessionId,
		UserId:        userId,
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

// getChatHistory retrieves the chat history for a session and converts it to Gemini content format
func (s *chatService) getChatHistory(ctx context.Context, chatSessionId, userId string) ([]*genai.Content, error) {
	s.logging.LogInfo(fmt.Sprintf("Retrieving chat history for session: %s", chatSessionId))

	chatDetails, err := s.chatRepository.FindChatSessionDetailByChatSessionIdAndUserId(chatSessionId, userId)
	if err != nil {
		s.logging.LogWarn(fmt.Sprintf("Failed to get chat history: %s", err.Error()))
		return []*genai.Content{}, nil // Return empty history if error, don't fail the whole request
	}

	var contents []*genai.Content

	// Convert chat history to Gemini content format
	for _, detail := range chatDetails {
		var role genai.Role
		if detail.Sender == models.ChatMessageSenderUser {
			role = genai.RoleUser
		} else {
			role = genai.RoleModel
		}

		content := genai.NewContentFromParts([]*genai.Part{
			genai.NewPartFromText(detail.Message),
		}, role)

		contents = append(contents, content)
	}

	s.logging.LogInfo(fmt.Sprintf("Retrieved %d messages from chat history", len(contents)))
	return contents, nil
}

func (s *chatService) SendChatMessage(ctx context.Context, req *models.ChatMessageRequest) (*responses.ChatMessageResponse, error) {
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

	err := s.chatRepository.InsertChatMessage(&models.ChatMessage{
		ChatMessageId: ulid.Make().String(),
		ChatSessionId: req.ChatSessionId,
		Message:       req.Message,
		Sender:        models.ChatMessageSenderUser,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		DeletedAt:     time.Time{},
	})
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to insert chat message: %s", err.Error()))
		return nil, fmt.Errorf("failed to insert chat message: %w", err)
	}

	var responseAi *responses.ChatMessageResponse

	// Handle different modes
	switch req.Mode {
	case models.ModeAgent:
		// Use RunAgent for agent mode
		s.logging.LogInfo("Using RunAgent for agent mode")

		// Get chat history for context
		chatDetails, err := s.chatRepository.FindChatSessionDetailByChatSessionIdAndUserId(req.ChatSessionId, req.UserId)
		if err != nil {
			s.logging.LogWarn(fmt.Sprintf("Failed to get chat history: %s", err.Error()))
		}

		// Build context message with chat history
		var messageWithContext string
		if len(chatDetails) > 0 {
			s.logging.LogInfo(fmt.Sprintf("Including %d previous messages for context", len(chatDetails)))
			contextStr := "\n\n--- CHAT HISTORY CONTEXT ---\n"
			for _, detail := range chatDetails {
				if detail.Sender == models.ChatMessageSenderUser {
					contextStr += fmt.Sprintf("User: %s\n", detail.Message)
				} else {
					contextStr += fmt.Sprintf("Assistant: %s\n", detail.Message)
				}
			}
			contextStr += "--- END CHAT HISTORY ---\n\n"
			messageWithContext = contextStr + "Current message: " + req.Message
		} else {
			messageWithContext = req.Message
		}

		response, err := s.geminiClient.RunAgent(ctx, messageWithContext, req.UserId)
		if err != nil {
			s.logging.LogError(fmt.Sprintf("Failed to run Gemini agent: %s", err.Error()))
			return nil, fmt.Errorf("failed to run Gemini agent: %w", err)
		}

		if err := s.logAIResponse(req.Message, response, req.UserId); err != nil {
			return nil, fmt.Errorf("failed to log AI response: %w", err)
		}

		input := openai.EmbeddingNewParamsInputUnion{
			OfString: param.NewOpt(req.Message),
		}
		embedding := s.openaiClient.CreateEmbedding(ctx, input)
		if embedding == nil {
			s.logging.LogError("Failed to create embedding: returned nil")
			return nil, fmt.Errorf("failed to create embedding")
		}

		err = s.chatRepository.InsertChatMessage(&models.ChatMessage{
			ChatMessageId: ulid.Make().String(),
			ChatSessionId: req.ChatSessionId,
			Message:       response.Response.(string),
			Sender:        models.ChatMessageSenderAssistant,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			DeletedAt:     time.Time{},
		})
		if err != nil {
			s.logging.LogError(fmt.Sprintf("Failed to insert AI chat message: %s", err.Error()))
			return nil, fmt.Errorf("failed to insert AI chat message: %w", err)
		}

		responseAi = &responses.ChatMessageResponse{
			ChatSessionId: req.ChatSessionId,
			ChatMessageId: ulid.Make().String(),
			Conversation: []*responses.Conversation{
				{
					Sender: models.ChatMessageSenderUser,
					Text:   req.Message,
				},
				{
					Sender: models.ChatMessageSenderAssistant,
					Text:   response.Response.(string),
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

	case models.ModeChat:
		fallthrough
	default:
		// Use regular Run for chat mode
		s.logging.LogInfo("Using Run for chat mode")

		// Get appropriate system prompt based on mode with user knowledge using RAG
		systemPrompt, err := s.getSystemPromptWithKnowledge(req.Mode, req.UserId, req.Message, ctx)
		if err != nil {
			s.logging.LogWarn(fmt.Sprintf("Failed to get enhanced system prompt: %s", err.Error()))
			systemPrompt = s.getSystemPromptByMode(req.Mode) // Fallback to base prompt
		}

		// Get chat history
		chatHistory, err := s.getChatHistory(ctx, req.ChatSessionId, req.UserId)
		if err != nil {
			s.logging.LogWarn(fmt.Sprintf("Failed to get chat history: %s", err.Error()))
			chatHistory = []*genai.Content{} // Use empty history if error
		}

		// Build message with system prompt, chat history, and current message
		message := []*genai.Content{
			genai.NewContentFromParts([]*genai.Part{
				genai.NewPartFromText(systemPrompt),
			}, genai.RoleModel),
		}

		// Add chat history
		message = append(message, chatHistory...)

		// Add current user message
		message = append(message, genai.NewContentFromParts([]*genai.Part{
			genai.NewPartFromText(req.Message),
		}, genai.RoleUser))

		response, err := s.geminiClient.Run(ctx, "gemini-2.5-flash", message)
		if err != nil {
			s.logging.LogError(fmt.Sprintf("Failed to run Gemini client: %s", err.Error()))
			return nil, fmt.Errorf("failed to run Gemini client: %w", err)
		}

		input := openai.EmbeddingNewParamsInputUnion{
			OfString: param.NewOpt(req.Message),
		}
		embedding := s.openaiClient.CreateEmbedding(ctx, input)
		if embedding == nil {
			s.logging.LogError("Failed to create embedding: returned nil")
			return nil, fmt.Errorf("failed to create embedding")
		}

		err = s.chatRepository.InsertChatMessage(&models.ChatMessage{
			ChatMessageId: ulid.Make().String(),
			ChatSessionId: req.ChatSessionId,
			Message:       response.Response.(string),
			Sender:        models.ChatMessageSenderAssistant,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			DeletedAt:     time.Time{},
		})
		if err != nil {
			s.logging.LogError(fmt.Sprintf("Failed to insert AI chat message: %s", err.Error()))
			return nil, fmt.Errorf("failed to insert AI chat message: %w", err)
		}

		responseAi = &responses.ChatMessageResponse{
			ChatSessionId: req.ChatSessionId,
			ChatMessageId: ulid.Make().String(),
			Conversation: []*responses.Conversation{
				{
					Sender: models.ChatMessageSenderUser,
					Text:   req.Message,
				},
				{
					Sender: models.ChatMessageSenderAssistant,
					Text:   response.Response.(string),
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}

	s.logging.LogInfo("Chat message processed successfully")
	return responseAi, nil
}

func (c *chatService) logAIResponse(responseString string, responseAi *responses.ResponseAI, userId string) error {
	messagePromptJSON, err := json.Marshal(responseString)
	if err != nil {
		c.logging.LogError(fmt.Sprintf("Failed to marshal message prompt: %v", err))
		return fmt.Errorf("failed to marshal message prompt: %w", err)
	}

	dateNow := time.Now()
	logMessage := &models.LogMessage{
		LogMessageId: ulid.Make().String(),
		UserId:       userId,
		Message:      string(messagePromptJSON),
		Response:     responseString,
		InputToken:   responseAi.InputToken,
		OutputToken:  responseAi.OutputToken,
		Topic:        "agent chat",
		Model:        "gemini-2.5-flash",
		CreatedAt:    dateNow,
		UpdatedAt:    dateNow,
	}

	err = c.logMessageService.InsertLogMessage(logMessage)
	if err != nil {
		c.logging.LogError(fmt.Sprintf("Failed to insert log message: %v", err))
		return fmt.Errorf("failed to insert log message: %w", err)
	}

	return nil
}

func (s *chatService) calculateCosineSimilarity(embedding1, embedding2 []float64) float64 {
	if len(embedding1) != len(embedding2) {
		return 0.0
	}

	var dotProduct, norm1, norm2 float64
	for i := 0; i < len(embedding1); i++ {
		dotProduct += embedding1[i] * embedding2[i]
		norm1 += embedding1[i] * embedding1[i]
		norm2 += embedding2[i] * embedding2[i]
	}

	if norm1 == 0.0 || norm2 == 0.0 {
		return 0.0
	}

	return dotProduct / (math.Sqrt(norm1) * math.Sqrt(norm2))
}

func (s *chatService) parseEmbedding(embeddingStr string) ([]float64, error) {
	// Remove brackets and split by comma
	embeddingStr = strings.Trim(embeddingStr, "[]")
	if embeddingStr == "" {
		return nil, fmt.Errorf("empty embedding string")
	}

	parts := strings.Split(embeddingStr, ",")
	embedding := make([]float64, len(parts))

	for i, part := range parts {
		val, err := strconv.ParseFloat(strings.TrimSpace(part), 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse embedding value: %w", err)
		}
		embedding[i] = val
	}

	return embedding, nil
}

func (s *chatService) createQueryEmbedding(ctx context.Context, query string) ([]float64, error) {
	input := openai.EmbeddingNewParamsInputUnion{
		OfString: param.NewOpt(query),
	}

	embedding := s.openaiClient.CreateEmbedding(ctx, input)

	if embedding == nil {
		return nil, fmt.Errorf("failed to create embedding")
	}

	return s.parseEmbedding(embedding.Embeddings)
}

func (s *chatService) gatherRelevantFinancialData(ctx context.Context, userId, query string) (*models.RelevantFinancialData, error) {
	s.logging.LogInfo(fmt.Sprintf("Gathering relevant financial data for user: %s", userId))

	// Create embedding for user query
	queryEmbedding, err := s.createQueryEmbedding(ctx, query)
	if err != nil {
		s.logging.LogWarn(fmt.Sprintf("Failed to create query embedding: %s", err.Error()))
		return nil, err
	}

	relevantData := &models.RelevantFinancialData{}

	// Get all user transactions
	transactionQuery := &requests.GetAllTransactionsQuery{
		Limit: 100, // Get more transactions for better RAG
	}

	transactionResponse, err := s.transactionService.GetAllTransactions(transactionQuery, userId)
	if err != nil {
		s.logging.LogWarn(fmt.Sprintf("Failed to get transactions: %s", err.Error()))
	} else if transactionResponse != nil {
		// Calculate similarity scores for transactions
		for i := range transactionResponse.Transactions {
			tx := &transactionResponse.Transactions[i]

			// Create searchable text from transaction
			searchText := fmt.Sprintf("%s %s %s", tx.Description, tx.Source, tx.Type)

			// Create embedding for transaction (if not already exists)
			var txEmbedding []float64
			if tx.DescriptionEmbedding != nil {
				if embeddingStr, ok := tx.DescriptionEmbedding.(string); ok {
					txEmbedding, err = s.parseEmbedding(embeddingStr)
					if err != nil {
						s.logging.LogWarn(fmt.Sprintf("Failed to parse transaction embedding: %s", err.Error()))
						continue
					}
				}
			} else {
				// Create new embedding for transaction
				txEmbedding, err = s.createQueryEmbedding(ctx, searchText)
				if err != nil {
					s.logging.LogWarn(fmt.Sprintf("Failed to create transaction embedding: %s", err.Error()))
					continue
				}
			}

			// Calculate similarity
			similarity := s.calculateCosineSimilarity(queryEmbedding, txEmbedding)

			// Only include transactions with similarity above threshold
			if similarity > 0.5 {
				relevantData.Transactions = append(relevantData.Transactions, models.TransactionWithScore{
					Transaction: tx,
					Score:       similarity,
				})
			}
		}
	}

	// Get all user receipts
	receipts, err := s.receiptService.GetReceiptsByUserId(userId)
	if err != nil {
		s.logging.LogWarn(fmt.Sprintf("Failed to get receipts: %s", err.Error()))
	} else {
		// Calculate similarity scores for receipts
		for _, receipt := range receipts {
			// Create searchable text from receipt
			searchText := fmt.Sprintf("%s %s", receipt.MerchantName, receipt.TransactionDate.Format("2006-01-02"))

			// Create embedding for receipt (if not already exists)
			var receiptEmbedding []float64
			if receipt.ExtractedReceiptEmbedding != nil {
				if embeddingStr, ok := receipt.ExtractedReceiptEmbedding.(string); ok {
					receiptEmbedding, err = s.parseEmbedding(embeddingStr)
					if err != nil {
						s.logging.LogWarn(fmt.Sprintf("Failed to parse receipt embedding: %s", err.Error()))
						continue
					}
				}
			} else {
				// Create new embedding for receipt
				receiptEmbedding, err = s.createQueryEmbedding(ctx, searchText)
				if err != nil {
					s.logging.LogWarn(fmt.Sprintf("Failed to create receipt embedding: %s", err.Error()))
					continue
				}
			}

			// Calculate similarity
			similarity := s.calculateCosineSimilarity(queryEmbedding, receiptEmbedding)

			// Only include receipts with similarity above threshold
			if similarity > 0.5 {
				relevantData.Receipts = append(relevantData.Receipts, models.ReceiptWithScore{
					Receipt: receipt,
					Score:   similarity,
				})
			}
		}
	}

	// Sort by similarity score (highest first)
	// Sort transactions
	for i := 0; i < len(relevantData.Transactions); i++ {
		for j := i + 1; j < len(relevantData.Transactions); j++ {
			if relevantData.Transactions[i].Score < relevantData.Transactions[j].Score {
				relevantData.Transactions[i], relevantData.Transactions[j] = relevantData.Transactions[j], relevantData.Transactions[i]
			}
		}
	}

	// Sort receipts
	for i := 0; i < len(relevantData.Receipts); i++ {
		for j := i + 1; j < len(relevantData.Receipts); j++ {
			if relevantData.Receipts[i].Score < relevantData.Receipts[j].Score {
				relevantData.Receipts[i], relevantData.Receipts[j] = relevantData.Receipts[j], relevantData.Receipts[i]
			}
		}
	}

	// Limit to top results for efficiency
	if len(relevantData.Transactions) > 15 {
		relevantData.Transactions = relevantData.Transactions[:15]
	}
	if len(relevantData.Receipts) > 10 {
		relevantData.Receipts = relevantData.Receipts[:10]
	}

	s.logging.LogInfo(fmt.Sprintf("Found %d relevant transactions and %d relevant receipts",
		len(relevantData.Transactions), len(relevantData.Receipts)))

	return relevantData, nil
}

func (s *chatService) gatherUserKnowledge(ctx context.Context, userId string) (*models.UserKnowledge, error) {
	s.logging.LogInfo(fmt.Sprintf("Gathering user knowledge for user: %s", userId))

	knowledge := &models.UserKnowledge{}

	// Get user transactions (limited to recent ones for context)
	transactionQuery := &requests.GetAllTransactionsQuery{
		Limit: 50, // Limit to recent 50 transactions
	}

	transactionResponse, err := s.transactionService.GetAllTransactions(transactionQuery, userId)
	if err != nil {
		s.logging.LogWarn(fmt.Sprintf("Failed to get transactions for user %s: %s", userId, err.Error()))
		// Don't fail completely, just log warning and continue
	} else if transactionResponse != nil {
		// Convert []models.Transaction to []*models.Transaction
		transactions := make([]*models.Transaction, len(transactionResponse.Transactions))
		for i := range transactionResponse.Transactions {
			transactions[i] = &transactionResponse.Transactions[i]
		}
		knowledge.Transactions = transactions
	}

	// Get user receipts (limited to recent ones for context)
	receipts, err := s.receiptService.GetReceiptsByUserId(userId)
	if err != nil {
		s.logging.LogWarn(fmt.Sprintf("Failed to get receipts for user %s: %s", userId, err.Error()))
		// Don't fail completely, just log warning and continue
	} else {
		// Limit receipts to recent 20 for context
		if len(receipts) > 20 {
			knowledge.Receipts = receipts[:20]
		} else {
			knowledge.Receipts = receipts
		}
	}

	s.logging.LogInfo(fmt.Sprintf("Gathered knowledge: %d transactions, %d receipts",
		len(knowledge.Transactions), len(knowledge.Receipts)))

	return knowledge, nil
}

func (s *chatService) buildKnowledgeContext(knowledge *models.UserKnowledge) string {
	if knowledge == nil {
		return ""
	}

	context := "\n\n--- USER'S FINANCIAL DATA CONTEXT ---\n"

	// Add transaction context
	if len(knowledge.Transactions) > 0 {
		context += fmt.Sprintf("\nRECENT TRANSACTIONS (%d):\n", len(knowledge.Transactions))
		for i, tx := range knowledge.Transactions {
			if i >= 10 { // Limit display to 10 most recent for prompt efficiency
				context += fmt.Sprintf("... and %d more transactions\n", len(knowledge.Transactions)-10)
				break
			}
			context += fmt.Sprintf("- %s: %s Rp %.0f (%s) - %s\n",
				tx.TransactionDate.Format("2006-01-02"),
				tx.Type,
				float64(tx.Amount), // Amount in Rupiah (no conversion needed)
				tx.Source,
				tx.Description)
		}
	}

	// Add receipt context
	if len(knowledge.Receipts) > 0 {
		context += fmt.Sprintf("\nRECENT RECEIPTS (%d):\n", len(knowledge.Receipts))
		for i, receipt := range knowledge.Receipts {
			if i >= 5 { // Limit display to 5 most recent for prompt efficiency
				context += fmt.Sprintf("... and %d more receipts\n", len(knowledge.Receipts)-5)
				break
			}
			context += fmt.Sprintf("- %s: %s - Total: Rp %.0f (Discount: Rp %.0f)\n",
				receipt.TransactionDate.Format("2006-01-02"),
				receipt.MerchantName,
				float64(receipt.TotalShopping), // Amount in Rupiah (no conversion needed)
				float64(receipt.TotalDiscount))
		}
	}

	context += "\n--- END OF FINANCIAL DATA CONTEXT ---\n\n"

	return context
}

func (s *chatService) buildRelevantKnowledgeContext(relevantData *models.RelevantFinancialData) string {
	if relevantData == nil {
		return ""
	}

	context := "\n\n--- RELEVANT FINANCIAL DATA CONTEXT (RAG) ---\n"

	// Add relevant transaction context
	if len(relevantData.Transactions) > 0 {
		context += fmt.Sprintf("\nMOST RELEVANT TRANSACTIONS (%d):\n", len(relevantData.Transactions))
		for i, txWithScore := range relevantData.Transactions {
			if i >= 10 { // Limit display to 10 most relevant for prompt efficiency
				context += fmt.Sprintf("... and %d more relevant transactions\n", len(relevantData.Transactions)-10)
				break
			}
			tx := txWithScore.Transaction
			context += fmt.Sprintf("- %s: %s Rp %.0f (%s) - %s (Relevance: %.2f)\n",
				tx.TransactionDate.Format("2006-01-02"),
				tx.Type,
				float64(tx.Amount), // Amount in Rupiah (no conversion needed)
				tx.Source,
				tx.Description,
				txWithScore.Score)
		}
	}

	// Add relevant receipt context
	if len(relevantData.Receipts) > 0 {
		context += fmt.Sprintf("\nMOST RELEVANT RECEIPTS (%d):\n", len(relevantData.Receipts))
		for i, receiptWithScore := range relevantData.Receipts {
			if i >= 5 { // Limit display to 5 most relevant for prompt efficiency
				context += fmt.Sprintf("... and %d more relevant receipts\n", len(relevantData.Receipts)-5)
				break
			}
			receipt := receiptWithScore.Receipt
			context += fmt.Sprintf("- %s: %s - Total: Rp %.0f (Discount: Rp %.0f) (Relevance: %.2f)\n",
				receipt.TransactionDate.Format("2006-01-02"),
				receipt.MerchantName,
				float64(receipt.TotalShopping), // Amount in Rupiah (no conversion needed)
				float64(receipt.TotalDiscount),
				receiptWithScore.Score)
		}
	}

	context += "\n--- END OF RELEVANT FINANCIAL DATA CONTEXT ---\n\n"

	return context
}

func (s *chatService) getSystemPromptWithKnowledge(mode models.Mode, userId string, userQuery string, ctx context.Context) (string, error) {
	basePrompt := s.getSystemPromptByMode(mode)

	// Only enhance Ask mode with user knowledge
	if mode != models.ModeChat {
		return basePrompt, nil
	}

	// Gather relevant user knowledge using RAG
	relevantData, err := s.gatherRelevantFinancialData(ctx, userId, userQuery)
	if err != nil {
		s.logging.LogWarn(fmt.Sprintf("Failed to gather relevant financial data: %s", err.Error()))
		return basePrompt, nil // Return base prompt if RAG fails
	}

	// Build relevant knowledge context
	knowledgeContext := s.buildRelevantKnowledgeContext(relevantData)

	// Check if we have relevant data
	hasRelevantData := len(relevantData.Transactions) > 0 || len(relevantData.Receipts) > 0

	if !hasRelevantData {
		// Fallback to basic knowledge if no relevant data found
		s.logging.LogInfo("No relevant financial data found, using basic knowledge gathering")
		basicKnowledge, err := s.gatherUserKnowledge(ctx, userId)
		if err == nil {
			knowledgeContext = s.buildKnowledgeContext(basicKnowledge)
		}
	}

	// Enhance prompt with knowledge
	var enhancedPrompt string
	if hasRelevantData {
		enhancedPrompt = basePrompt + knowledgeContext +
			"Use the above RELEVANT financial data to provide more personalized and accurate responses. " +
			"The data has been selected based on semantic similarity to the user's question. " +
			"Reference specific transactions, receipts, or patterns when relevant to the user's question. " +
			"The relevance scores indicate how closely each item matches the user's query."
	} else {
		enhancedPrompt = basePrompt + knowledgeContext +
			"Use the above financial data to provide more personalized and accurate responses. " +
			"Reference specific transactions, receipts, or patterns when relevant to the user's question."
	}

	return enhancedPrompt, nil
}
