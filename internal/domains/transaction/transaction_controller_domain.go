package transaction

import "github.com/gofiber/fiber/v2"

type TransactionController interface {
	GetAllTransactions(ctx *fiber.Ctx) error
	CreateTransaction(ctx *fiber.Ctx) error
	GetDetailedTransaction(ctx *fiber.Ctx) error
	UpdateTransaction(ctx *fiber.Ctx) error
	DeleteTransaction(ctx *fiber.Ctx) error
	OverviewTransactions(ctx *fiber.Ctx) error
}
