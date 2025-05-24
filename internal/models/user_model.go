package models

import "time"

type User struct {
	UserId    string    `json:"user_id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RegisterUser struct {
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

type FindUserById struct {
	UserId   string `json:"user_id" validate:"required"`
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type UpdateUserRequest struct {
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}
