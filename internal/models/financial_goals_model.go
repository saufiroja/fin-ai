package models

import "time"

type FinancialGoal struct {
	FinancialGoalId string    `json:"financial_goal_id"`
	UserId          string    `json:"user_id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	TargetAmount    float64   `json:"target_amount"`
	CurrentAmount   float64   `json:"current_amount"`
	TargetDate      time.Time `json:"target_date"`
	Status          string    `json:"status"` // e.g., "active", "achieved", "cancelled"
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
