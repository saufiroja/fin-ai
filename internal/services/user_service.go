package services

import (
	"errors"
	"fmt"

	"github.com/saufiroja/fin-ai/internal/interfaces"
	"github.com/saufiroja/fin-ai/internal/models"
	logging "github.com/saufiroja/fin-ai/pkg/loggings"
)

type userService struct {
	UserRepository interfaces.UserRepositoryInterface
	logging        logging.Logger
}

func NewUserService(userRepository interfaces.UserRepositoryInterface, logger logging.Logger) interfaces.UserServiceInterface {
	return &userService{
		UserRepository: userRepository,
		logging:        logger,
	}
}

func (s *userService) UpdateUserById(userId string, req *models.UpdateUserRequest) error {
	s.logging.LogInfo(fmt.Sprintf("Updating user with ID: %s", userId))
	_, err := s.UserRepository.FindUserById(userId)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("User with ID %s not found: %v", userId, err))
		return errors.New("user not found")
	}

	err = s.UserRepository.UpdateUserById(userId, req)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to update user with ID %s: %v", userId, err))
		return errors.New("failed to update user")
	}

	s.logging.LogInfo(fmt.Sprintf("User with ID %s updated successfully", userId))
	return nil
}

func (s *userService) DeleteUserById(userId string) error {
	s.logging.LogInfo(fmt.Sprintf("Deleting user with ID: %s", userId))
	_, err := s.UserRepository.FindUserById(userId)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("User with ID %s not found: %v", userId, err))
		return errors.New("user not found")
	}

	err = s.UserRepository.DeleteUserById(userId)
	if err != nil {
		s.logging.LogError(fmt.Sprintf("Failed to delete user with ID %s: %v", userId, err))
		return errors.New("failed to delete user")
	}

	s.logging.LogInfo(fmt.Sprintf("User with ID %s deleted successfully", userId))
	return nil
}

func (s *userService) GetMe(userId string) (*models.FindUserById, error) {
	s.logging.LogInfo("Getting user information")

	user, err := s.UserRepository.FindUserById(userId)
	if err != nil {
		s.logging.LogError("Failed to get user info: " + err.Error())
		return nil, errors.New("failed to get user information")
	}

	s.logging.LogInfo("User information retrieved successfully")
	return user, nil
}
