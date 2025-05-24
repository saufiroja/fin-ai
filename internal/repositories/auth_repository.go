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

func (r *authRepository) FindUserByEmail(email string) (*models.User, error) {
	db := r.DB.Connection()
	query := `SELECT user_id, full_name, email, password FROM users WHERE email = $1`
	row := db.QueryRow(query, email)

	var user models.User
	err := row.Scan(&user.UserId, &user.FullName, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *authRepository) FindUserById(userId string) (*models.FindUserById, error) {
	db := r.DB.Connection()
	query := `SELECT user_id, full_name, email FROM users WHERE user_id = $1`
	row := db.QueryRow(query, userId)

	var user models.FindUserById
	err := row.Scan(&user.UserId, &user.FullName, &user.Email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
