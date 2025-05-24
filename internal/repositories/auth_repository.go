package repositories

import (
	"github.com/saufiroja/fin-ai/internal/interfaces"
	"github.com/saufiroja/fin-ai/internal/models"
	"github.com/saufiroja/fin-ai/pkg/databases"
)

type userRepository struct {
	DB databases.PostgresManager
}

func NewUserRepository(db databases.PostgresManager) interfaces.UserRepositoryInterface {
	return &userRepository{
		DB: db,
	}
}

func (r *userRepository) InsertUser(req *models.User) error {
	db := r.DB.Connection()
	query := `INSERT INTO users (user_id, full_name, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.Exec(query, req.UserId, req.FullName, req.Email, req.Password, req.CreatedAt, req.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) FindUserByEmail(email string) (*models.User, error) {
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

func (r *userRepository) FindUserById(userId string) (*models.FindUserById, error) {
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

func (r *userRepository) UpdateUserById(userId string, req *models.UpdateUserRequest) error {
	db := r.DB.Connection()
	query := `UPDATE users SET 
			full_name = $1,
			email = $2, 
			updated_at = NOW() 
			WHERE user_id = $3`
	_, err := db.Exec(query, req.FullName, req.Email, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) DeleteUserById(userId string) error {
	db := r.DB.Connection()
	query := `DELETE FROM users WHERE user_id = $1`
	_, err := db.Exec(query, userId)
	if err != nil {
		return err
	}

	return nil
}
