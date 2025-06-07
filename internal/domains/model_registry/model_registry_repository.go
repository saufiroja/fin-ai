package model_registry

import "github.com/saufiroja/fin-ai/internal/models"

type ModelRegistryStorer interface {
	FindAllModels() ([]*models.ModelRegistry, error)
	FindModelById(modelId string) (*models.ModelRegistry, error)
}
