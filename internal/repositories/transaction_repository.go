package repositories

import (
	"github.com/saufiroja/fin-ai/internal/contracts/requests"
	"github.com/saufiroja/fin-ai/internal/contracts/responses"
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

func (t *transactionRepository) GetAllTransactions(req *requests.GetAllTransactionsQuery, userId string) ([]models.Transaction, error) {
	db := t.DB.Connection()

	query := `
        SELECT 
            transaction_id, user_id, category_id, type, amount, 
            description, description_embedding, source, transaction_date, 
            ai_category_confidence, is_auto_categorized, created_at, updated_at,
            confirmed, discount
        FROM transactions
        WHERE ($1 = '' OR category_id = $1)
        AND ($2 = '' OR LOWER(description) LIKE LOWER('%' || $2 || '%'))
        AND (NULLIF($6, '') IS NULL OR NULLIF($7, '') IS NULL OR 
             transaction_date BETWEEN 
             ($6::date + INTERVAL '0 hours')::timestamp AND 
             ($7::date + INTERVAL '23 hours 59 minutes 59 seconds')::timestamp)
		AND user_id = $5
        ORDER BY transaction_date DESC
        LIMIT $3 OFFSET $4`

	rows, err := db.Query(query, req.CategoryId, req.Search, req.Limit, req.Offset, userId, req.StartDate, req.EndDate)
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
			&transaction.Confirmed,
			&transaction.Discount,
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
    updated_at,
	confirmed,
	discount,
	payment_method
    )
    VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
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
		transaction.Confirmed,
		transaction.Discount,
		transaction.PaymentMethod,
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
        updated_at,
        confirmed,
        discount,
		payment_method
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
		&transaction.Confirmed,
		&transaction.Discount,
		&transaction.PaymentMethod,
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
        updated_at = $11,
        confirmed = $12,
        discount = $13,
		payment_method = $14
    WHERE transaction_id = $15
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
		transaction.Confirmed,
		transaction.Discount,
		transaction.PaymentMethod,
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

func (t *transactionRepository) CountAllTransactions(req *requests.GetAllTransactionsQuery, userId string) (int64, error) {
	db := t.DB.Connection()

	query := `
	SELECT COUNT(*)
	FROM transactions
	WHERE ($1 = '' OR category_id = $1)
	AND ($2 = '' OR LOWER(description) LIKE LOWER('%' || $2 || '%'))
	AND user_id = $3
	AND (NULLIF($4, '') IS NULL OR NULLIF($5, '') IS NULL 
	OR transaction_date BETWEEN ($4::date + INTERVAL '0 hours')::timestamp AND 
	($5::date + INTERVAL '23 hours 59 minutes 59 seconds')::timestamp)
	`

	var count int64
	err := db.QueryRow(query, req.CategoryId, req.Search, userId, req.StartDate, req.EndDate).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (t *transactionRepository) GetTransactionsStats(userId string, req *requests.OverviewTransactionsQuery) (*responses.OverviewTransactions, error) {
	db := t.DB.Connection()

	query := `
        SELECT 
            COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) AS total_income,
            COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) AS total_expense
        FROM transactions
		WHERE user_id = $1
		AND ($4 = '' OR category_id = $4)
        AND (NULLIF($2, '') IS NULL OR NULLIF($3, '') IS NULL 
		OR transaction_date BETWEEN ($2::date + INTERVAL '0 hours')::timestamp AND 
		($3::date + INTERVAL '23 hours 59 minutes 59 seconds')::timestamp)
    `

	row := db.QueryRow(query, userId, req.StartDate, req.EndDate, req.CategoryId)

	stats := &responses.OverviewTransactions{}
	err := row.Scan(
		&stats.TotalIncome,
		&stats.TotalExpense,
	)

	if err != nil {
		return nil, err
	}

	return stats, nil
}
