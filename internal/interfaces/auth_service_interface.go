package interfaces

import "github.com/saufiroja/fin-ai/internal/models"

type AuthServiceInterface interface {
	RegisterUser(req *models.User) error
}
