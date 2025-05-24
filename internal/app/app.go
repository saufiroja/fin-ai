package app

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/saufiroja/fin-ai/config"
	authController "github.com/saufiroja/fin-ai/internal/controllers/auth"
	user "github.com/saufiroja/fin-ai/internal/controllers/user"
	"github.com/saufiroja/fin-ai/internal/middleware"
	"github.com/saufiroja/fin-ai/internal/repositories"
	"github.com/saufiroja/fin-ai/internal/services"
	"github.com/saufiroja/fin-ai/internal/utils"
	"github.com/saufiroja/fin-ai/pkg/databases"
	logging "github.com/saufiroja/fin-ai/pkg/loggings"
)

type App struct {
	*fiber.App
}

func NewApp() *App {
	return &App{
		App: fiber.New(),
	}
}

func (a *App) Start() {
	logger := logging.NewLogrusAdapter()
	conf := config.NewAppConfig(logger)
	postgresInstance := databases.NewPostgres(conf, logger)
	authMiddleware := middleware.Authorization(conf)
	// llm.NewOpenAI(conf)
	defer func() {
		if err := postgresInstance.CloseConnection(); err != nil {
			logger.LogError(fmt.Sprintf("failed to close postgres connection: %v", err))
		}
	}()

	// health check
	globalApi := a.Group("/api/v1")
	globalApi.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "ok",
		})
	})

	validator := utils.NewValidator()
	tokenGenerator := utils.NewJWTTokenGenerator(conf)
	userRepository := repositories.NewUserRepository(postgresInstance)

	authService := services.NewAuthService(userRepository, logger, tokenGenerator, conf)
	userService := services.NewUserService(userRepository, logger)

	authController := authController.NewAuthController(authService, validator)
	userController := user.NewUserController(userService)

	auth := globalApi.Group("/auth")
	auth.Post("/register", authController.RegisterUser)
	auth.Post("/login", authController.LoginUser)
	auth.Post("/logout", authController.LogoutUser)
	auth.Post("/refresh-token", authController.RefreshToken)

	user := globalApi.Group("/user")
	user.Get("/me", authMiddleware, userController.GetMe)
	user.Put("/:user_id", authMiddleware, userController.UpdateUserById)
	user.Delete("/:user_id", authMiddleware, userController.DeleteUserById)

	if err := a.Listen(fmt.Sprintf("localhost:%s", conf.Http.Port)); err != nil {
		logger.LogPanic(err.Error())
	}
}
