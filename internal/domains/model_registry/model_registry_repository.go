package model_registry

import "github.com/saufiroja/fin-ai/internal/models"

type ModelRegistryRepository interface {
	FindAllModels() ([]*models.ModelRegistry, error)
	FindModelById(modelId string) (*models.ModelRegistry, error)
}
