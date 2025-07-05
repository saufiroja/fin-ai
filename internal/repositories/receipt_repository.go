package repositories

import (
	"fmt"

	"github.com/saufiroja/fin-ai/internal/contracts/requests"
	"github.com/saufiroja/fin-ai/internal/domains/receipt"
	"github.com/saufiroja/fin-ai/internal/models"
	"github.com/saufiroja/fin-ai/pkg/databases"
)

type receiptRepository struct {
	DB databases.PostgresManager
}

func NewReceiptRepository(db databases.PostgresManager) receipt.ReceiptStorer {
	return &receiptRepository{
		DB: db,
	}
}
func (r *receiptRepository) InsertReceipt(receipt *models.Receipt) error {
	db := r.DB.Connection()

	query := `
    INSERT INTO receipts (
    receipt_id, 
    user_id,
    merchant_name,
    sub_total,
    total_discount,
    total_shopping,
    metadata,
    extracted_receipt, 
    extracted_receipt_embedding, 
    confirmed,
    transaction_date, 
    created_at, 
    updated_at
    ) 
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err := db.Exec(query, receipt.ReceiptId, receipt.UserId, receipt.MerchantName, receipt.SubTotal, receipt.TotalDiscount, receipt.TotalShopping, receipt.MetaData, receipt.ExtractedReceipt, receipt.ExtractedReceiptEmbedding, receipt.Confirmed, receipt.TransactionDate, receipt.CreatedAt, receipt.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *receiptRepository) InsertReceiptItem(receiptItem *models.ReceiptItem) error {
	db := r.DB.Connection()

	query := `
    INSERT INTO receipt_items (
    receipt_item_id, 
    receipt_id, 
    item_name, 
    item_quantity, 
    item_price, 
    item_price_total, 
    item_discount, 
    created_at, 
    updated_at,
	category_id,
	ai_category_confidence
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW(), $8, $9)`

	_, err := db.Exec(query, receiptItem.ReceiptItemId, receiptItem.ReceiptId, receiptItem.ItemName, receiptItem.ItemQuantity, receiptItem.ItemPrice, receiptItem.ItemPriceTotal, receiptItem.ItemDiscount, receiptItem.CategoryId, receiptItem.AiCategoryConfidence)
	if err != nil {
		return err
	}

	return nil
}

func (r *receiptRepository) GetReceiptsByUserId(userId string) ([]*models.Receipt, error) {
	db := r.DB.Connection()

	query := `
    SELECT 
        receipt_id, 
        user_id, 
        merchant_name, 
        sub_total, 
        total_discount, 
        total_shopping, 
        metadata, 
        extracted_receipt, 
        extracted_receipt_embedding, 
        confirmed, 
        transaction_date, 
        created_at, 
        updated_at
    FROM receipts
    WHERE user_id = $1`

	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var receipts []*models.Receipt
	for rows.Next() {
		var receipt models.Receipt
		if err := rows.Scan(&receipt.ReceiptId, &receipt.UserId, &receipt.MerchantName, &receipt.SubTotal, &receipt.TotalDiscount, &receipt.TotalShopping, &receipt.MetaData, &receipt.ExtractedReceipt, &receipt.ExtractedReceiptEmbedding, &receipt.Confirmed, &receipt.TransactionDate, &receipt.CreatedAt, &receipt.UpdatedAt); err != nil {
			return nil, err
		}
		receipts = append(receipts, &receipt)
	}

	return receipts, nil
}

func (r *receiptRepository) GetAllReceiptsByUserId(userId string, req *requests.GetAllReceiptsQuery) ([]*models.Receipt, error) {
	db := r.DB.Connection()
	fmt.Println("GetAllReceiptsByUserId called with userId:", userId, "and request:", req)

	// Build the ORDER BY clause safely
	orderBy := "created_at DESC" // default
	if req.SortBy != "" && req.SortOrder != "" {
		// Validate sortBy to prevent SQL injection
		validSortColumns := map[string]bool{
			"created_at":       true,
			"updated_at":       true,
			"transaction_date": true,
			"merchant_name":    true,
			"total_shopping":   true,
			"sub_total":        true,
		}

		validSortOrders := map[string]bool{
			"ASC":  true,
			"DESC": true,
		}

		if validSortColumns[req.SortBy] && validSortOrders[req.SortOrder] {
			orderBy = req.SortBy + " " + req.SortOrder
		}
	}

	query := fmt.Sprintf(`
    SELECT 
        receipt_id, 
        user_id, 
        merchant_name, 
        sub_total, 
        total_discount, 
        total_shopping, 
        metadata, 
        extracted_receipt, 
        extracted_receipt_embedding, 
        confirmed, 
        transaction_date, 
        created_at, 
        updated_at
    FROM receipts
    WHERE user_id = $1
	AND (merchant_name ILIKE '%%' || $2 || '%%' OR $2 = '')
	ORDER BY %s
	LIMIT $3 OFFSET $4`, orderBy)

	rows, err := db.Query(query, userId, req.Search, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var receipts []*models.Receipt
	for rows.Next() {
		var receipt models.Receipt
		if err := rows.Scan(&receipt.ReceiptId, &receipt.UserId, &receipt.MerchantName, &receipt.SubTotal, &receipt.TotalDiscount, &receipt.TotalShopping, &receipt.MetaData, &receipt.ExtractedReceipt, &receipt.ExtractedReceiptEmbedding, &receipt.Confirmed, &receipt.TransactionDate, &receipt.CreatedAt, &receipt.UpdatedAt); err != nil {
			return nil, err
		}
		receipts = append(receipts, &receipt)
	}

	return receipts, nil
}

func (r *receiptRepository) GetDetailReceiptUserById(userId string, receiptId string) (*models.Receipt, error) {
	db := r.DB.Connection()

	query := `
    SELECT 
        r.receipt_id,
        r.user_id,
        r.merchant_name,
        r.sub_total,
        r.total_discount,
        r.total_shopping,
        r.confirmed,
        r.transaction_date,
        r.created_at,
        r.updated_at        
    FROM receipts r
    WHERE r.user_id = $1 AND r.receipt_id = $2`

	row := db.QueryRow(query, userId, receiptId)

	var receipt models.Receipt
	if err := row.Scan(
		&receipt.ReceiptId,
		&receipt.UserId,
		&receipt.MerchantName,
		&receipt.SubTotal,
		&receipt.TotalDiscount,
		&receipt.TotalShopping,
		&receipt.Confirmed,
		&receipt.TransactionDate,
		&receipt.CreatedAt,
		&receipt.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &receipt, nil
}

func (r *receiptRepository) GetReceiptItemsByReceiptId(receiptId string) ([]*models.ReceiptItem, error) {
	db := r.DB.Connection()
	query := `
    SELECT 
        receipt_item_id, 
        receipt_id, 
        item_name, 
        item_quantity, 
        item_price, 
        item_price_total, 
        item_discount, 
        created_at, 
        updated_at,
		category_id,
		ai_category_confidence
    FROM receipt_items
    WHERE receipt_id = $1`

	rows, err := db.Query(query, receiptId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.ReceiptItem
	for rows.Next() {
		var item models.ReceiptItem
		if err := rows.Scan(&item.ReceiptItemId, &item.ReceiptId, &item.ItemName, &item.ItemQuantity, &item.ItemPrice, &item.ItemPriceTotal, &item.ItemDiscount, &item.CreatedAt, &item.UpdatedAt, &item.CategoryId, &item.AiCategoryConfidence); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	return items, nil
}

func (r *receiptRepository) UpdateReceiptConfirmed(receiptId string, confirmed bool) error {
	db := r.DB.Connection()

	query := `
	UPDATE receipts
	SET confirmed = $1, updated_at = NOW()
	WHERE receipt_id = $2`

	_, err := db.Exec(query, confirmed, receiptId)
	if err != nil {
		return err
	}

	return nil
}

func (r *receiptRepository) CountReceiptsByUserId(userId string, req *requests.GetAllReceiptsQuery) (int64, error) {
	db := r.DB.Connection()

	query := `
	SELECT COUNT(*)
	FROM receipts
	WHERE user_id = $1
	AND (merchant_name ILIKE '%' || $2 || '%' OR $2 = '')`

	var count int64
	err := db.QueryRow(query, userId, req.Search).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
