package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/saufiroja/fin-ai/internal/interfaces"
	"github.com/saufiroja/fin-ai/internal/models"
	logging "github.com/saufiroja/fin-ai/pkg/loggings"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	repoUser interfaces.AuthRepositoryInterface
	logging  logging.Logger
}

func NewAuthService(repo interfaces.AuthRepositoryInterface, logging logging.Logger) interfaces.AuthServiceInterface {
	return &authService{
		repoUser: repo,
		logging:  logging,
	}
}

func (s *authService) RegisterUser(req *models.User) error {
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
