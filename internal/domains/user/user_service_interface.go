package user

import (
	"github.com/saufiroja/fin-ai/internal/contracts/requests"
	"github.com/saufiroja/fin-ai/internal/contracts/responses"
)

type UserService interface {
	UpdateUserById(userId string, req *requests.UpdateUserRequest) error
	DeleteUserById(userId string) error
	GetMe(userId string) (*responses.FindUserById, error)
}
