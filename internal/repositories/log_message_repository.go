package repositories

import (
	"github.com/saufiroja/fin-ai/internal/domains/log_message"
	"github.com/saufiroja/fin-ai/internal/models"
	"github.com/saufiroja/fin-ai/pkg/databases"
)

type logMessageRepository struct {
	DB databases.PostgresManager
}

func NewLogMessageRepository(db databases.PostgresManager) log_message.LogMessageRepository {
	return &logMessageRepository{
		DB: db,
	}
}

func (r *logMessageRepository) InsertLogMessage(logMessage *models.LogMessageModel) error {
	db := r.DB.Connection()

	query := `INSERT INTO log_messages (log_messages_id, user_id, message, response, input_token, output_token, topic, model, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())`

	_, err := db.Exec(query, logMessage.LogMessageId, logMessage.UserId, logMessage.Message, logMessage.Response,
		logMessage.InputToken, logMessage.OutputToken, logMessage.Topic, logMessage.Model)
	if err != nil {
		return err
	}

	return nil
}
