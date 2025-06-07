package repositories

import (
	"github.com/saufiroja/fin-ai/internal/domains"
	"github.com/saufiroja/fin-ai/internal/models"
	"github.com/saufiroja/fin-ai/pkg/databases"
)

type modelRegistryRepository struct {
	DB databases.PostgresManager
}

func NewModelRegistryRepository(db databases.PostgresManager) domains.ModelRegistryRepositoryInterface {
	return &modelRegistryRepository{
		DB: db,
	}
}

func (r *modelRegistryRepository) FindAllModels() ([]*models.ModelRegistry, error) {
	db := r.DB.Connection()

	var modelsRegistry []*models.ModelRegistry
	query := `SELECT model_registry_id, name, created_at, updated_at FROM model_registries ORDER BY created_at DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var model models.ModelRegistry
		err := rows.Scan(&model.ModelRegistryId, &model.Name, &model.CreatedAt, &model.UpdatedAt)
		if err != nil {
			return nil, err
		}
		modelsRegistry = append(modelsRegistry, &model)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return modelsRegistry, nil
}

func (r *modelRegistryRepository) FindModelById(modelId string) (*models.ModelRegistry, error) {
	db := r.DB.Connection()

	var model models.ModelRegistry
	query := `SELECT model_registry_id, name, created_at, updated_at FROM model_registries WHERE model_registry_id = $1`
	err := db.QueryRow(query, modelId).Scan(&model.ModelRegistryId, &model.Name, &model.CreatedAt, &model.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &model, nil
}
