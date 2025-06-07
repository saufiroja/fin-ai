package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saufiroja/fin-ai/internal/contracts/requests"
	"github.com/saufiroja/fin-ai/internal/contracts/responses"
	"github.com/saufiroja/fin-ai/internal/domains/transaction"
)

type transactionController struct {
	transactionService transaction.TransactionManager
}

func NewTransactionController(transactionService transaction.TransactionManager) transaction.TransactionController {
	return &transactionController{
		transactionService: transactionService,
	}
}

func (t *transactionController) GetAllTransactions(ctx *fiber.Ctx) error {
	transactionQuery := &requests.GetAllTransactionsQuery{
		Limit:    10, // Default limit
		Offset:   0,  // Default offset
		Category: "",
		Search:   "",
	}

	transactions, err := t.transactionService.GetAllTransactions(transactionQuery)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to retrieve transactions",
		})
	}

	return ctx.JSON(transactions)
}

// CreateTransaction implements transaction.TransactionController.
func (t *transactionController) CreateTransaction(ctx *fiber.Ctx) error {
	panic("unimplemented")
}

// DeleteTransaction implements transaction.TransactionController.
func (t *transactionController) DeleteTransaction(ctx *fiber.Ctx) error {
	panic("unimplemented")
}

// GetDetailedTransaction implements transaction.TransactionController.
func (t *transactionController) GetDetailedTransaction(ctx *fiber.Ctx) error {
	panic("unimplemented")
}

// GetTransactionsStats implements transaction.TransactionController.
func (t *transactionController) GetTransactionsStats(ctx *fiber.Ctx) error {
	panic("unimplemented")
}

// UpdateTransaction implements transaction.TransactionController.
func (t *transactionController) UpdateTransaction(ctx *fiber.Ctx) error {
	panic("unimplemented")
}
