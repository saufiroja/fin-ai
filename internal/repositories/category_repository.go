package repositories

import (
	"github.com/saufiroja/fin-ai/internal/domains/categories"
	"github.com/saufiroja/fin-ai/internal/models"
	"github.com/saufiroja/fin-ai/pkg/databases"
)

type categoryRepository struct {
	DB databases.PostgresManager
}

func NewCategoryRepository(db databases.PostgresManager) categories.CategoryStorer {
	return &categoryRepository{
		DB: db,
	}
}

func (c *categoryRepository) InsertCategory(category *models.Category) error {
	db := c.DB.Connection()

	query := `
    INSERT INTO categories (category_id, name, name_embedding, type, created_at, updated_at)
    VALUES ($1, $2, $3, $4, NOW(), NOW())`
	_, err := db.Exec(query, category.CategoryId, category.Name, category.NameEmbedding, category.Type)
	if err != nil {
		return err
	}

	return nil
}
