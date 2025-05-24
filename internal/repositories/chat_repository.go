package repositories

import (
	"github.com/saufiroja/fin-ai/internal/interfaces"
	"github.com/saufiroja/fin-ai/internal/models"
	"github.com/saufiroja/fin-ai/pkg/databases"
)

type chatRepository struct {
	DB databases.PostgresManager
}

func NewChatRepository(db databases.PostgresManager) interfaces.ChatRepositoryInterface {
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
	query := `SELECT chat_session_id, user_id, title FROM chat_sessions WHERE user_id = $1 ORDER BY updated_at DESC`
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
