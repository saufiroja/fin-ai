package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saufiroja/fin-ai/internal/contracts/requests"
	"github.com/saufiroja/fin-ai/internal/contracts/responses"
	"github.com/saufiroja/fin-ai/internal/domains/receipt"
)

type receiptController struct {
	receiptService receipt.ReceiptManager
}

func NewReceiptController(receiptService receipt.ReceiptManager) receipt.ReceiptController {
	return &receiptController{
		receiptService: receiptService,
	}
}

func (r *receiptController) UploadReceipt(c *fiber.Ctx) error {
	file, err := c.FormFile("receipt")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Failed to get file from request",
		})
	}

	userId := c.Locals("user_id").(string)

	receipt, err := r.receiptService.UploadReceipt(file, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to upload receipt",
		})
	}

	return c.Status(fiber.StatusOK).JSON(responses.Response{
		Status:  fiber.StatusOK,
		Message: "Receipt uploaded successfully",
		Data:    receipt,
	})
}

func (r *receiptController) GetReceiptsByUserId(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	query := &requests.GetAllReceiptsQuery{
		Limit:     10, // Default limit
		Offset:    1,  // Default offset
		Search:    "",
		SortBy:    "created_at", // Default sort by
		SortOrder: "DESC",       // Default sort order
	}
	if err := c.QueryParser(query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid query parameters",
		})
	}

	receipts, err := r.receiptService.GetAllReceiptsByUserId(userId, query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to get receipts",
		})
	}

	return c.Status(fiber.StatusOK).JSON(responses.Response{
		Status:  fiber.StatusOK,
		Message: "Receipts retrieved successfully",
		Data:    receipts.Receipts,
		Pagination: &responses.Pagination{
			Total:       receipts.Total,
			TotalPages:  receipts.TotalPages,
			CurrentPage: receipts.CurrentPage,
		},
	})
}

func (r *receiptController) GetDetailReceiptUserById(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	receiptId := c.Params("receipt_id")

	detail, err := r.receiptService.GetDetailReceiptUserById(userId, receiptId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to get receipt details",
		})
	}

	return c.Status(fiber.StatusOK).JSON(responses.Response{
		Status:  fiber.StatusOK,
		Message: "Receipt details retrieved successfully",
		Data:    detail,
	})
}

func (r *receiptController) UpdateReceiptConfirmed(c *fiber.Ctx) error {
	receiptId := c.Params("receipt_id")
	confirmed := c.Query("confirmed") == "true"
	userId := c.Locals("user_id").(string)

	err := r.receiptService.UpdateReceiptConfirmed(userId, receiptId, confirmed)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to update receipt confirmation",
		})
	}

	return c.Status(fiber.StatusOK).JSON(responses.Response{
		Status:  fiber.StatusOK,
		Message: "Receipt confirmation updated successfully",
	})
}
