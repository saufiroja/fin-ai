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

	chatGroup.Post("/session/:user_id", r.container.Dependencies.AuthMiddleware, r.container.Controllers.Chat.CreateChatSession)
	chatGroup.Get("/session/:user_id", r.container.Dependencies.AuthMiddleware, r.container.Controllers.Chat.FindAllChatSessions)
	chatGroup.Put("/session-rename", r.container.Dependencies.AuthMiddleware, r.container.Controllers.Chat.RenameChatSession)
	chatGroup.Delete("/session/:chat_session_id/:user_id", r.container.Dependencies.AuthMiddleware, r.container.Controllers.Chat.DeleteChatSession)
	chatGroup.Post("/send", r.container.Dependencies.AuthMiddleware, r.container.Controllers.Chat.SendChatMessage)
	chatGroup.Get("/session-detail/:chat_session_id/:user_id", r.container.Dependencies.AuthMiddleware, r.container.Controllers.Chat.GetChatSessionDetail)
}
