package models

import "time"

type User struct {
	UserId        string    `json:"user_id"`
	FullName      string    `json:"full_name"`
	Email         string    `json:"email"`
	Password      string    `json:"password"`
	AiPreferences any       `json:"ai_preferences,omitempty"` // Use 'any' for flexible AI preferences storage
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
