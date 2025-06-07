package repositories

import (
	"github.com/saufiroja/fin-ai/internal/contracts/requests"
	"github.com/saufiroja/fin-ai/internal/domains/transaction"
	"github.com/saufiroja/fin-ai/internal/models"
	"github.com/saufiroja/fin-ai/pkg/databases"
)

type transactionRepository struct {
	DB databases.PostgresManager
}

func NewTransactionRepository(db databases.PostgresManager) transaction.TransactionStorer {
	return &transactionRepository{
		DB: db,
	}
}

func (t *transactionRepository) GetAllTransactions(req *requests.GetAllTransactionsQuery) ([]models.Transaction, error) {
	db := t.DB.Connection()

	query := `
        SELECT 
            transaction_id, user_id, category_id, type, amount, 
            description, description_embedding, source, transaction_date, 
            ai_category_confidence, is_auto_categorized, created_at, updated_at
        FROM transactions
        WHERE ($1 = '' OR category_id = $1)
        AND ($2 = '' OR LOWER(description) LIKE LOWER('%' || $2 || '%'))
        ORDER BY transaction_date DESC
        LIMIT $3 OFFSET $4`

	rows, err := db.Query(query, req.Category, req.Search, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		transaction := models.Transaction{}
		err := rows.Scan(
			&transaction.TransactionId,
			&transaction.UserId,
			&transaction.CategoryId,
			&transaction.Type,
			&transaction.Amount,
			&transaction.Description,
			&transaction.DescriptionEmbedding,
			&transaction.Source,
			&transaction.TransactionDate,
			&transaction.AiCategoryConfidence,
			&transaction.IsAutoCategorized,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (t *transactionRepository) InsertTransaction(transaction *models.Transaction) error {
	db := t.DB.Connection()

	query := `
    INSERT INTO transactions (
    transaction_id, 
    user_id, 
    category_id,
    type,
    description,
    description_embedding,
    amount,
    source,
    transaction_date,
    ai_category_confidence,
    is_auto_categorized,
    created_at,
    updated_at
    )
    VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
	)
`
	_, err := db.Exec(query,
		transaction.TransactionId,
		transaction.UserId,
		transaction.CategoryId,
		transaction.Type,
		transaction.Description,
		transaction.DescriptionEmbedding,
		transaction.Amount,
		transaction.Source,
		transaction.TransactionDate,
		transaction.AiCategoryConfidence,
		transaction.IsAutoCategorized,
		transaction.CreatedAt,
		transaction.UpdatedAt,
	)

	return err
}

func (t *transactionRepository) GetTransactionByID(id string) (*models.Transaction, error) {
	db := t.DB.Connection()

	query := `
    SELECT 
        transaction_id, 
        user_id, 
        category_id, 
        type, 
        amount, 
        description, 
        description_embedding, 
        source, 
        transaction_date, 
        ai_category_confidence, 
        is_auto_categorized, 
        created_at, 
        updated_at
    FROM transactions
    WHERE transaction_id = $1
`

	row := db.QueryRow(query, id)

	transaction := &models.Transaction{}
	err := row.Scan(
		&transaction.TransactionId,
		&transaction.UserId,
		&transaction.CategoryId,
		&transaction.Type,
		&transaction.Amount,
		&transaction.Description,
		&transaction.DescriptionEmbedding,
		&transaction.Source,
		&transaction.TransactionDate,
		&transaction.AiCategoryConfidence,
		&transaction.IsAutoCategorized,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *transactionRepository) UpdateTransaction(transaction *models.Transaction) error {
	db := t.DB.Connection()

	query := `
    UPDATE transactions
    SET 
        user_id = $1,
        category_id = $2,
        type = $3,
        amount = $4,
        description = $5,
        description_embedding = $6,
        source = $7,
        transaction_date = $8,
        ai_category_confidence = $9,
        is_auto_categorized = $10,
        updated_at = $11
    WHERE transaction_id = $12
`

	_, err := db.Exec(query,
		transaction.UserId,
		transaction.CategoryId,
		transaction.Type,
		transaction.Amount,
		transaction.Description,
		transaction.DescriptionEmbedding,
		transaction.Source,
		transaction.TransactionDate,
		transaction.AiCategoryConfidence,
		transaction.IsAutoCategorized,
		transaction.UpdatedAt,
		transaction.TransactionId,
	)

	return err
}

func (t *transactionRepository) DeleteTransaction(id string) error {
	db := t.DB.Connection()

	query := `
    DELETE FROM transactions
    WHERE transaction_id = $1
`

	_, err := db.Exec(query, id)

	return err
}

func (t *transactionRepository) CountAllTransactions(req *requests.GetAllTransactionsQuery) (int64, error) {
	db := t.DB.Connection()

	query := `
        SELECT COUNT(*)
        FROM transactions
        WHERE ($1 = '' OR category_id = $1)
        AND ($2 = '' OR LOWER(description) LIKE LOWER('%' || $2 || '%'))`

	var count int64
	err := db.QueryRow(query, req.Category, req.Search).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
