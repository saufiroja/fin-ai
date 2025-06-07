package domains

import "github.com/saufiroja/fin-ai/internal/models"

type UserRepositoryInterface interface {
	InsertUser(req *models.User) error
	FindUserByEmail(email string) (*models.User, error)
	FindUserById(userId string) (*models.FindUserById, error)
	UpdateUserById(userId string, req *models.UpdateUserRequest) error
	DeleteUserById(userId string) error
}
