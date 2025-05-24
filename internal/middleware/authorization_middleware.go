package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/saufiroja/fin-ai/config"
	"github.com/saufiroja/fin-ai/internal/models"
)

func Authorization(conf *config.AppConfig) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Ambil token dari cookie atau header
		tokenString := ctx.Cookies("access_token")
		secret := conf.Jwt.Secret
		if tokenString == "" {
			// Coba ambil dari Authorization header
			authHeader := ctx.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenString == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(models.Response{
				Status:  fiber.StatusUnauthorized,
				Message: "Unauthorized: No token provided",
			})
		}

		// Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan algoritma cocok
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return ctx.Status(fiber.StatusUnauthorized).JSON(models.Response{
				Status:  fiber.StatusUnauthorized,
				Message: "Unauthorized: Invalid token",
			})
		}

		// Simpan klaim ke context untuk digunakan di handler
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx.Locals("user_id", claims["user_id"])
			ctx.Locals("email", claims["email"])
			ctx.Locals("full_name", claims["full_name"])
			ctx.Locals("roles", claims["roles"])
		}

		return ctx.Next()
	}
}
