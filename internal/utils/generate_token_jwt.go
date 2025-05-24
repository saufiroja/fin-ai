package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/saufiroja/fin-ai/config"
	"github.com/saufiroja/fin-ai/internal/models"
)

// TokenGenerator defines methods to generate and validate JWT tokens
type TokenGenerator interface {
	GenerateAccessToken(userId, fullName, email string) (*models.JwtGenerator, error)
	GenerateRefreshToken(userId, fullName, email string) (*models.JwtGenerator, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
}

// JWTTokenGenerator implements TokenGenerator interface
type JWTTokenGenerator struct {
	conf *config.AppConfig
}

// NewJWTTokenGenerator creates a new JWT token generator
func NewJWTTokenGenerator(conf *config.AppConfig) TokenGenerator {
	return &JWTTokenGenerator{
		conf: conf,
	}
}

// GenerateAccessToken generates a JWT access token with custom claims
func (g *JWTTokenGenerator) GenerateAccessToken(userId, fullName, email string) (*models.JwtGenerator, error) {
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	// Create standard claims
	claims := jwt.MapClaims{
		"user_id":   userId,
		"full_name": fullName,
		"email":     email,
		"iat":       now.Unix(),
		"exp":       expiresAt.Unix(),
		"iss":       "fin-ai",
		"typ":       "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(g.conf.Jwt.Secret))
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT access token: %w", err)
	}

	return &models.JwtGenerator{
		Token:     tokenString,
		ExpiredAt: expiresAt,
	}, nil
}

// GenerateRefreshToken generates a JWT refresh token
func (g *JWTTokenGenerator) GenerateRefreshToken(userId, fullName, email string) (*models.JwtGenerator, error) {
	now := time.Now()
	expiresAt := now.Add(30 * 24 * time.Hour) // 30 days

	claims := jwt.MapClaims{
		"user_id":   userId,
		"full_name": fullName,
		"email":     email,
		"iat":       now.Unix(),
		"exp":       expiresAt.Unix(),
		"iss":       "fin-ai",
		"typ":       "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(g.conf.Jwt.Secret))
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT refresh token: %w", err)
	}

	return &models.JwtGenerator{
		Token:     tokenString,
		ExpiredAt: expiresAt,
	}, nil
}

// ValidateToken validates a JWT token
func (g *JWTTokenGenerator) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(g.conf.Jwt.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return token, nil
}
