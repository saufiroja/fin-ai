package controllers

import (
	"github.com/gofiber/fiber/v2"
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

	err = r.receiptService.UploadReceipt(file, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to upload receipt",
		})
	}

	return c.Status(fiber.StatusOK).JSON(responses.Response{
		Status:  fiber.StatusOK,
		Message: "Receipt uploaded successfully",
	})
}

func (r *receiptController) GetReceiptsByUserId(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)

	receipts, err := r.receiptService.GetReceiptsByUserId(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to get receipts",
		})
	}

	return c.Status(fiber.StatusOK).JSON(responses.Response{
		Status:  fiber.StatusOK,
		Message: "Receipts retrieved successfully",
		Data:    receipts,
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

	err := r.receiptService.UpdateReceiptConfirmed(receiptId, confirmed)
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
