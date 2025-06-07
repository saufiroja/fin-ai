package repositories

import (
	"github.com/saufiroja/fin-ai/internal/domains"
	"github.com/saufiroja/fin-ai/internal/models"
	"github.com/saufiroja/fin-ai/pkg/databases"
)

type chatRepository struct {
	DB databases.PostgresManager
}

func NewChatRepository(db databases.PostgresManager) domains.ChatRepositoryInterface {
	return &chatRepository{
		DB: db,
	}
}

func (r *chatRepository) InsertChatSession(chatSession *models.ChatSession) error {
	db := r.DB.Connection()
	query := `INSERT INTO chat_sessions (chat_session_id, user_id, title, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(query, chatSession.ChatSessionId, chatSession.UserId, "New Chat", chatSession.CreatedAt, chatSession.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *chatRepository) FindAllChatSessions(userId string) ([]*models.ChatSession, error) {
	db := r.DB.Connection()
	query := `SELECT chat_session_id, user_id, title FROM chat_sessions WHERE user_id = $1 AND deleted_at IS NULL ORDER BY updated_at DESC`
	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chatSessions []*models.ChatSession
	for rows.Next() {
		var chatSession models.ChatSession
		err := rows.Scan(&chatSession.ChatSessionId, &chatSession.UserId, &chatSession.Title)
		if err != nil {
			return nil, err
		}
		chatSessions = append(chatSessions, &chatSession)
	}

	return chatSessions, nil
}

func (r *chatRepository) RenameChatSession(chatSession *models.ChatSession) error {
	db := r.DB.Connection()
	query := `UPDATE chat_sessions SET title = $1, updated_at = $2 WHERE chat_session_id = $3 AND user_id = $4`
	_, err := db.Exec(query, chatSession.Title, chatSession.UpdatedAt, chatSession.ChatSessionId, chatSession.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (r *chatRepository) DeleteChatSession(chatSessionId, userId string) error {
	db := r.DB.Connection()
	query := `UPDATE chat_sessions SET deleted_at = NOW() WHERE chat_session_id = $1 AND user_id = $2`
	_, err := db.Exec(query, chatSessionId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (r *chatRepository) FindChatSessionByChatSessionIdAndUserId(chatSessionId, userId string) (*models.ChatSession, error) {
	db := r.DB.Connection()
	query := `SELECT chat_session_id, user_id, title FROM chat_sessions WHERE chat_session_id = $1 AND user_id = $2`
	row := db.QueryRow(query, chatSessionId, userId)

	var chatSession models.ChatSession
	err := row.Scan(&chatSession.ChatSessionId, &chatSession.UserId, &chatSession.Title)
	if err != nil {
		return nil, err
	}

	return &chatSession, nil
}

func (r *chatRepository) InsertChatMessage(chatMessage *models.ChatMessage) error {
	db := r.DB.Connection()
	query := `INSERT INTO chat_messages (chat_message_id, chat_session_id, message, sender, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.Exec(query, chatMessage.ChatMessageId, chatMessage.ChatSessionId, chatMessage.Message,
		chatMessage.Sender, chatMessage.CreatedAt, chatMessage.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *chatRepository) FindChatMessagesByChatSessionId(chatSessionId string) ([]*models.ChatMessage, error) {
	db := r.DB.Connection()
	query := `SELECT chat_message_id, chat_session_id, message, sender, created_at FROM chat_messages WHERE chat_session_id = $1 ORDER BY created_at ASC`
	rows, err := db.Query(query, chatSessionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chatMessages []*models.ChatMessage
	for rows.Next() {
		var chatMessage models.ChatMessage
		err := rows.Scan(&chatMessage.ChatMessageId, &chatMessage.ChatSessionId, &chatMessage.Message,
			&chatMessage.Sender, &chatMessage.CreatedAt)
		if err != nil {
			return nil, err
		}
		chatMessages = append(chatMessages, &chatMessage)
	}

	return chatMessages, nil
}

func (r *chatRepository) UpdateChatSessionTitle(chatSession *models.ChatSession) error {
	db := r.DB.Connection()
	query := `UPDATE chat_sessions SET title = $1, updated_at = $2 WHERE chat_session_id = $3 AND user_id = $4`
	_, err := db.Exec(query, chatSession.Title, chatSession.UpdatedAt, chatSession.ChatSessionId, chatSession.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (r *chatRepository) FindChatSessionDetailByChatSessionIdAndUserId(chatSessionId, userId string) ([]*models.ChatSessionDetail, error) {
	db := r.DB.Connection()
	query := `SELECT cm.chat_message_id, cm.chat_session_id, cm.message, cm.sender, cm.created_at
	FROM chat_sessions cs
	JOIN chat_messages cm ON cs.chat_session_id = cm.chat_session_id
	WHERE cs.chat_session_id = $1 AND cs.user_id = $2 AND cs.deleted_at IS NULL
	ORDER BY cm.created_at ASC`
	rows, err := db.Query(query, chatSessionId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chatMessages []*models.ChatSessionDetail
	for rows.Next() {
		var chatMessage models.ChatSessionDetail
		err := rows.Scan(&chatMessage.ChatMessageId, &chatMessage.ChatSessionId, &chatMessage.Message,
			&chatMessage.Sender, &chatMessage.CreatedAt)
		if err != nil {
			return nil, err
		}
		chatMessages = append(chatMessages, &chatMessage)
	}

	return chatMessages, nil
}
