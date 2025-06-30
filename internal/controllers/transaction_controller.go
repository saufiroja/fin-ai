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
	userId := ctx.Locals("user_id").(string) // Assuming user_id is set in context
	transactionQuery := &requests.GetAllTransactionsQuery{
		Limit:     10, // Default limit
		Offset:    1,  // Default offset
		Category:  "",
		Search:    "",
		StartDate: "",
		EndDate:   "",
	}
	if err := ctx.QueryParser(transactionQuery); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid query parameters",
		})
	}

	transactions, err := t.transactionService.GetAllTransactions(transactionQuery, userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to retrieve transactions",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.Response{
		Status:  fiber.StatusOK,
		Message: "Transactions retrieved successfully",
		Data:    transactions.Transactions,
		Pagination: &responses.Pagination{
			Total:       transactions.Total,
			CurrentPage: transactions.CurrentPage,
			TotalPages:  transactions.TotalPages,
		},
	})
}

func (t *transactionController) CreateTransaction(ctx *fiber.Ctx) error {
	req := &requests.TransactionRequest{
		UserId: ctx.Locals("user_id").(string), // Assuming user_id is set in context
	}

	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	if err := t.transactionService.InsertTransaction(req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to create transaction",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(responses.Response{
		Status:  fiber.StatusCreated,
		Message: "Transaction created successfully",
	})
}

func (t *transactionController) DeleteTransaction(ctx *fiber.Ctx) error {
	transactionId := ctx.Params("transaction_id")
	if transactionId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Transaction ID is required",
		})
	}

	if err := t.transactionService.DeleteTransaction(transactionId); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to delete transaction",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.Response{
		Status:  fiber.StatusOK,
		Message: "Transaction deleted successfully",
	})
}

func (t *transactionController) GetDetailedTransaction(ctx *fiber.Ctx) error {
	transactionId := ctx.Params("transaction_id")
	if transactionId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Transaction ID is required",
		})
	}

	transaction, err := t.transactionService.GetDetailedTransaction(transactionId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to retrieve transaction details",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.Response{
		Status:  fiber.StatusOK,
		Message: "Transaction details retrieved successfully",
		Data:    transaction,
	})
}

func (t *transactionController) UpdateTransaction(ctx *fiber.Ctx) error {
	transactionId := ctx.Params("transaction_id")
	if transactionId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Transaction ID is required",
		})
	}

	req := &requests.UpdateTransactionRequest{
		UserId: ctx.Locals("user_id").(string), // Assuming user_id is set in context
	}

	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	if err := t.transactionService.UpdateTransaction(transactionId, req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to update transaction",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.Response{
		Status:  fiber.StatusOK,
		Message: "Transaction updated successfully",
	})
}

func (t *transactionController) OverviewTransactions(ctx *fiber.Ctx) error {
	userId := ctx.Locals("user_id").(string) // Assuming user_id is set in context
	transactionQuery := &requests.OverviewTransactionsQuery{
		StartDate: "",
		EndDate:   "",
	}
	if err := ctx.QueryParser(transactionQuery); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid query parameters",
		})
	}

	overview, err := t.transactionService.OverviewTransactions(userId, transactionQuery)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to retrieve overview transactions",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.Response{
		Status:  fiber.StatusOK,
		Message: "Overview transactions retrieved successfully",
		Data:    overview,
	})
}
