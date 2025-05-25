package models

import "time"

type ModelRegistry struct {
	ModelRegistryId string     `json:"model_registry_id"`
	Name            string     `json:"name"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}
