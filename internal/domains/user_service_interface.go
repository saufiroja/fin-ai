package domains

import "github.com/saufiroja/fin-ai/internal/models"

type UserServiceInterface interface {
	UpdateUserById(userId string, req *models.UpdateUserRequest) error
	DeleteUserById(userId string) error
	GetMe(userId string) (*models.FindUserById, error)
}
