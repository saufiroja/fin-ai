package interfaces

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saufiroja/fin-ai/internal/models"
)

type AuthServiceInterface interface {
	RegisterUser(req *models.RegisterUser) error
	LoginUser(req *models.LoginUser, ctx *fiber.Ctx) (*models.LoginResponse, error)
	LogoutUser(ctx *fiber.Ctx) error
	ValidateRefreshToken(token string) (*models.JwtGenerator, error)
	GetMe(userId string) (*models.FindUserById, error)
}
