package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saufiroja/fin-ai/internal/contracts/requests"
	"github.com/saufiroja/fin-ai/internal/contracts/responses"
	"github.com/saufiroja/fin-ai/internal/models"
)

type AuthService interface {
	RegisterUser(req *requests.RegisterUser) error
	LoginUser(req *requests.LoginUser, ctx *fiber.Ctx) (*responses.LoginResponse, error)
	LogoutUser(ctx *fiber.Ctx) error
	ValidateRefreshToken(token string) (*models.JwtGenerator, error)
}
