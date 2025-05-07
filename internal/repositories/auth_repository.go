package repositories

import (
	"github.com/saufiroja/fin-ai/internal/interfaces"
	"github.com/saufiroja/fin-ai/internal/models"
	"github.com/saufiroja/fin-ai/pkg/databases"
)

type authRepository struct {
	DB databases.PostgresManager
}

func NewAuthRepository(db databases.PostgresManager) interfaces.AuthRepositoryInterface {
	return &authRepository{
		DB: db,
	}
}

func (r *authRepository) InsertUser(req *models.User) error {
	db := r.DB.Connection()
	query := `INSERT INTO users (user_id, full_name, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.Exec(query, req.UserId, req.FullName, req.Email, req.Password, req.CreatedAt, req.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}
