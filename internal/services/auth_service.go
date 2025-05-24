package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/saufiroja/fin-ai/internal/interfaces"
	"github.com/saufiroja/fin-ai/internal/models"
	"github.com/saufiroja/fin-ai/internal/utils"
	logging "github.com/saufiroja/fin-ai/pkg/loggings"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	repoUser       interfaces.AuthRepositoryInterface
	logging        logging.Logger
	tokenGenerator utils.TokenGenerator
}

func NewAuthService(repo interfaces.AuthRepositoryInterface, logging logging.Logger, tokenGenerator utils.TokenGenerator) interfaces.AuthServiceInterface {
	return &authService{
		repoUser:       repo,
		logging:        logging,
		tokenGenerator: tokenGenerator,
	}
}

func (s *authService) RegisterUser(req *models.RegisterUser) error {
	s.logging.LogInfo("Registering user")

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Error hashing password: %v", err))
		return errors.New("failed to hash password")
	}

	user := &models.User{
		UserId:    ulid.Make().String(),
		FullName:  req.FullName,
		Email:     req.Email,
		Password:  string(hash),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err = s.repoUser.InsertUser(user)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Error inserting user into database: %v", err))
		return errors.New("failed to register user")
	}

	s.logging.LogInfo("User registered successfully")
	return nil
}

func (s *authService) LoginUser(req *models.LoginUser) (*models.LoginResponse, error) {
	s.logging.LogInfo("Logging in user")

	user, err := s.repoUser.FindUserByEmail(req.Email)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Error finding user by email: %v", err))
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Error comparing password: %v", err))
		return nil, errors.New("invalid email or password")
	}

	accessToken, err := s.tokenGenerator.GenerateAccessToken(user.UserId, user.FullName, user.Email)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Error generating access token: %v", err))
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := s.tokenGenerator.GenerateRefreshToken(user.UserId, user.FullName, user.Email)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Error generating refresh token: %v", err))
		return nil, errors.New("failed to generate refresh token")
	}

	res := &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	s.logging.LogInfo("User logged in successfully")
	return res, nil
}
