package domains

import "github.com/saufiroja/fin-ai/internal/models"

type ModelRegistryRepositoryInterface interface {
	FindAllModels() ([]*models.ModelRegistry, error)
	FindModelById(modelId string) (*models.ModelRegistry, error)
}
