package services

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
	"github.com/saufiroja/fin-ai/config"
	"github.com/saufiroja/fin-ai/internal/contracts/requests"
	"github.com/saufiroja/fin-ai/internal/contracts/responses"
	"github.com/saufiroja/fin-ai/internal/domains"
	"github.com/saufiroja/fin-ai/internal/models"
	"github.com/saufiroja/fin-ai/internal/utils"
	logging "github.com/saufiroja/fin-ai/pkg/loggings"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	repoUser       domains.UserRepositoryInterface
	logging        logging.Logger
	tokenGenerator utils.TokenGenerator
	conf           *config.AppConfig
}

func NewAuthService(
	repo domains.UserRepositoryInterface,
	logging logging.Logger,
	tokenGenerator utils.TokenGenerator,
	conf *config.AppConfig,
) domains.AuthServiceInterface {
	return &authService{
		repoUser:       repo,
		logging:        logging,
		tokenGenerator: tokenGenerator,
		conf:           conf,
	}
}

func (s *authService) RegisterUser(req *requests.RegisterUser) error {
	s.logging.LogInfo("Registering user")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logging.LogError("Failed to hash password: " + err.Error())
		return errors.New("failed to hash password")
	}

	user := &models.User{
		UserId:    ulid.Make().String(),
		FullName:  req.FullName,
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.repoUser.InsertUser(user); err != nil {
		s.logging.LogError("Failed to insert user into database: " + err.Error())
		return errors.New("failed to register user")
	}

	s.logging.LogInfo("User registered successfully")
	return nil
}

func (s *authService) LoginUser(req *requests.LoginUser, ctx *fiber.Ctx) (*responses.LoginResponse, error) {
	s.logging.LogInfo("Logging in user")

	user, err := s.repoUser.FindUserByEmail(req.Email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		s.logging.LogError("Invalid email or password")
		return nil, errors.New("invalid email or password")
	}

	accessToken, err := s.tokenGenerator.GenerateAccessToken(user.UserId, user.FullName, user.Email)
	if err != nil {
		s.logging.LogError("Failed to generate access token: " + err.Error())
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := s.tokenGenerator.GenerateRefreshToken(user.UserId, user.FullName, user.Email)
	if err != nil {
		s.logging.LogError("Failed to generate refresh token: " + err.Error())
		return nil, errors.New("failed to generate refresh token")
	}

	res := &responses.LoginResponse{
		AccessToken:           accessToken.Token,
		RefreshToken:          refreshToken.Token,
		AccessTokenExpiresAt:  accessToken.ExpiredAt,
		RefreshTokenExpiresAt: refreshToken.ExpiredAt,
	}

	s.setAuthCookies(ctx, res)

	s.logging.LogInfo("User logged in successfully")
	return res, nil
}

func (s *authService) LogoutUser(ctx *fiber.Ctx) error {
	s.logging.LogInfo("Logging out user")
	s.clearAuthCookies(ctx)
	s.logging.LogInfo("User logged out successfully")
	return nil
}

func (s *authService) ValidateRefreshToken(token string) (*models.JwtGenerator, error) {
	s.logging.LogInfo("Validating refresh token")

	claims, err := s.tokenGenerator.ValidateToken(token)
	if err != nil {
		s.logging.LogError("Invalid refresh token: " + err.Error())
		return nil, errors.New("invalid refresh token")
	}

	mapClaims, ok := claims.Claims.(jwt.MapClaims)
	if !ok {
		s.logging.LogError("Failed to parse token claims")
		return nil, errors.New("invalid token claims")
	}

	userId := mapClaims["user_id"].(string)
	fullName := mapClaims["full_name"].(string)
	email := mapClaims["email"].(string)

	newAccessToken, err := s.tokenGenerator.GenerateAccessToken(userId, fullName, email)
	if err != nil {
		s.logging.LogError("Failed to generate new access token: " + err.Error())
		return nil, errors.New("failed to generate new access token")
	}

	s.logging.LogInfo("Refresh token validated successfully")
	return newAccessToken, nil
}

func (s *authService) setAuthCookies(ctx *fiber.Ctx, res *responses.LoginResponse) {
	domain := "localhost" // or use s.conf.AppDomain if configurable

	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    res.AccessToken,
		Domain:   domain,
		SameSite: "disabled",
		MaxAge:   int(time.Until(res.AccessTokenExpiresAt).Seconds()),
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    res.RefreshToken,
		Domain:   domain,
		SameSite: "disabled",
		MaxAge:   int(time.Until(res.RefreshTokenExpiresAt).Seconds()),
	})
}

func (s *authService) clearAuthCookies(ctx *fiber.Ctx) {
	domain := "localhost"

	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Domain:   domain,
		SameSite: "disabled",
		MaxAge:   -1,
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Domain:   domain,
		SameSite: "disabled",
		MaxAge:   -1,
	})
}
