package repositories

import (
	"github.com/saufiroja/fin-ai/internal/contracts/requests"
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

func (c *categoryRepository) FindAllCategories(req *requests.GetAllCategoryQuery) ([]models.Category, error) {
	db := c.DB.Connection()

	query := `
	SELECT category_id, name, name_embedding, type, created_at, updated_at
	FROM categories
	WHERE ($1::text IS NULL OR name ILIKE '%' || $1 || '%')
	ORDER BY created_at DESC
	LIMIT $2 OFFSET $3`

	rows, err := db.Query(query, req.Search, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categoriesList []models.Category
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.CategoryId, &category.Name, &category.NameEmbedding, &category.Type, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, err
		}
		categoriesList = append(categoriesList, category)
	}

	return categoriesList, nil
}

func (c *categoryRepository) CountCategories(req *requests.GetAllCategoryQuery) (int64, error) {
	db := c.DB.Connection()

	query := `
	SELECT COUNT(*)
	FROM categories
	WHERE ($1::text IS NULL OR name ILIKE '%' || $1 || '%')`

	var count int64
	err := db.QueryRow(query, req.Search).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
