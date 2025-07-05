package app

import (
	"github.com/gofiber/fiber/v2"
)

type Routes struct {
	app       *fiber.App
	container *Container
}

func NewRoutes(app *fiber.App, container *Container) *Routes {
	return &Routes{
		app:       app,
		container: container,
	}
}

func (r *Routes) Setup() {
	r.setupHealthCheck()
	r.setupAuthRoutes()
	r.setupUserRoutes()
	r.setupChatRoutes()
	r.setupTransactionRoutes()
	r.setupCategoryRoutes()
	r.setupReceiptRoutes()
}

func (r *Routes) setupHealthCheck() {
	globalApi := r.app.Group("/api/v1")
	globalApi.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "ok",
		})
	})
}

func (r *Routes) setupAuthRoutes() {
	globalApi := r.app.Group("/api/v1")
	authGroup := globalApi.Group("/auth")

	authGroup.Post("/register", r.container.Controllers.Auth.RegisterUser)
	authGroup.Post("/login", r.container.Controllers.Auth.LoginUser)
	authGroup.Post("/logout", r.container.Controllers.Auth.LogoutUser)
	authGroup.Post("/refresh-token", r.container.Controllers.Auth.RefreshToken)
}

func (r *Routes) setupUserRoutes() {
	globalApi := r.app.Group("/api/v1")
	userGroup := globalApi.Group("/user")

	userGroup.Get("/me", r.container.Dependencies.AuthMiddleware, r.container.Controllers.User.GetMe)
	userGroup.Put("/:user_id", r.container.Dependencies.AuthMiddleware, r.container.Controllers.User.UpdateUserById)
	userGroup.Delete("/:user_id", r.container.Dependencies.AuthMiddleware, r.container.Controllers.User.DeleteUserById)
}

func (r *Routes) setupChatRoutes() {
	globalApi := r.app.Group("/api/v1")
	chatGroup := globalApi.Group("/chat")

	chatGroup.Post("/sessions", r.container.Dependencies.AuthMiddleware, r.container.Controllers.Chat.CreateChatSession)
	chatGroup.Get("/sessions", r.container.Dependencies.AuthMiddleware, r.container.Controllers.Chat.FindAllChatSessions)
	chatGroup.Put("/sessions/rename/:chat_session_id", r.container.Dependencies.AuthMiddleware, r.container.Controllers.Chat.RenameChatSession)
	chatGroup.Delete("/sessions/:chat_session_id", r.container.Dependencies.AuthMiddleware, r.container.Controllers.Chat.DeleteChatSession)
	chatGroup.Post("/sessions/send", r.container.Dependencies.AuthMiddleware, r.container.Controllers.Chat.SendChatMessage)
	chatGroup.Get("/sessions/:chat_session_id", r.container.Dependencies.AuthMiddleware, r.container.Controllers.Chat.GetChatSessionDetail)
}

func (r *Routes) setupTransactionRoutes() {
	globalApi := r.app.Group("/api/v1")
	transactionGroup := globalApi.Group("/transactions")

	transactionGroup.Post("/",
		r.container.Dependencies.AuthMiddleware,
		r.container.Controllers.Transaction.CreateTransaction)
	transactionGroup.Get("/",
		r.container.Dependencies.AuthMiddleware,
		r.container.Controllers.Transaction.GetAllTransactions)
	transactionGroup.Get("/overviews",
		r.container.Dependencies.AuthMiddleware,
		r.container.Controllers.Transaction.OverviewTransactions)
	transactionGroup.Get("/:transaction_id",
		r.container.Dependencies.AuthMiddleware,
		r.container.Controllers.Transaction.GetDetailedTransaction)
	transactionGroup.Put("/:transaction_id",
		r.container.Dependencies.AuthMiddleware,
		r.container.Controllers.Transaction.UpdateTransaction)
	transactionGroup.Delete("/:transaction_id",
		r.container.Dependencies.AuthMiddleware,
		r.container.Controllers.Transaction.DeleteTransaction)
}

func (r *Routes) setupCategoryRoutes() {
	globalApi := r.app.Group("/api/v1")
	categoryGroup := globalApi.Group("/categories")

	categoryGroup.Post("/",
		r.container.Dependencies.AuthMiddleware,
		r.container.Controllers.Category.CreateCategory)
	categoryGroup.Get("/",
		r.container.Dependencies.AuthMiddleware,
		r.container.Controllers.Category.GetAllCategories)
	categoryGroup.Put("/:category_id",
		r.container.Dependencies.AuthMiddleware,
		r.container.Controllers.Category.UpdateCategoryById)
	categoryGroup.Delete("/:category_id",
		r.container.Dependencies.AuthMiddleware,
		r.container.Controllers.Category.DeleteCategoryById)
}

func (r *Routes) setupReceiptRoutes() {
	globalApi := r.app.Group("/api/v1")
	receiptGroup := globalApi.Group("/receipts")

	receiptGroup.Post("/upload",
		r.container.Dependencies.AuthMiddleware,
		r.container.Controllers.Receipt.UploadReceipt)
	receiptGroup.Get("/user",
		r.container.Dependencies.AuthMiddleware,
		r.container.Controllers.Receipt.GetReceiptsByUserId)
	receiptGroup.Get("/detail/user/:receipt_id",
		r.container.Dependencies.AuthMiddleware,
		r.container.Controllers.Receipt.GetDetailReceiptUserById)
	receiptGroup.Put("/confirm/:receipt_id",
		r.container.Dependencies.AuthMiddleware,
		r.container.Controllers.Receipt.UpdateReceiptConfirmed)
}
