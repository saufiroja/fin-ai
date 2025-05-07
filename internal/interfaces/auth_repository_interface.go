package interfaces

import "github.com/saufiroja/fin-ai/internal/models"

type AuthRepositoryInterface interface {
	InsertUser(req *models.User) error
}
