package user

import (
	"github.com/saufiroja/fin-ai/internal/contracts/requests"
	"github.com/saufiroja/fin-ai/internal/contracts/responses"
	"github.com/saufiroja/fin-ai/internal/models"
)

type UserStorer interface {
	InsertUser(req *models.User) error
	FindUserByEmail(email string) (*models.User, error)
	FindUserById(userId string) (*responses.FindUserById, error)
	UpdateUserById(userId string, req *requests.UpdateUserRequest) error
	DeleteUserById(userId string) error
}
