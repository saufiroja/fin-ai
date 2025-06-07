package services

import (
	"errors"
	"fmt"

	"github.com/saufiroja/fin-ai/internal/contracts/requests"
	"github.com/saufiroja/fin-ai/internal/contracts/responses"
	"github.com/saufiroja/fin-ai/internal/domains"
	logging "github.com/saufiroja/fin-ai/pkg/loggings"
)

type userService struct {
	UserRepository domains.UserRepositoryInterface
	logging        logging.Logger
}

func NewUserService(userRepository domains.UserRepositoryInterface, logger logging.Logger) domains.UserServiceInterface {
	return &userService{
		UserRepository: userRepository,
		logging:        logger,
	}
}

func (s *userService) UpdateUserById(userId string, req *requests.UpdateUserRequest) error {
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

func (s *userService) GetMe(userId string) (*responses.FindUserById, error) {
	s.logging.LogInfo("Getting user information")

	user, err := s.UserRepository.FindUserById(userId)
	if err != nil {
		s.logging.LogError("Failed to get user info: " + err.Error())
		return nil, errors.New("failed to get user information")
	}

	s.logging.LogInfo("User information retrieved successfully")
	return user, nil
}
