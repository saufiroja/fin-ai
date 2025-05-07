package app

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/saufiroja/fin-ai/config"
	authController "github.com/saufiroja/fin-ai/internal/controllers/auth"
	"github.com/saufiroja/fin-ai/internal/repositories"
	"github.com/saufiroja/fin-ai/internal/services"
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
	// llm.NewOpenAI(conf)
	defer func() {
		if err := postgresInstance.CloseConnection(); err != nil {
			logger.LogError(fmt.Sprintf("failed to close postgres connection: %v", err))
		}
	}()

	userRepository := repositories.NewAuthRepository(postgresInstance)
	userService := services.NewAuthService(userRepository, logger)
	userController := authController.NewAuthController(userService)

	a.Post("/register", userController.RegisterUser)

	if err := a.Listen(fmt.Sprintf(":%s", conf.Http.Port)); err != nil {
		logger.LogPanic(err.Error())
	}
}
