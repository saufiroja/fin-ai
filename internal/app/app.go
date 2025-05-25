package app

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/saufiroja/fin-ai/config"
	authController "github.com/saufiroja/fin-ai/internal/controllers/auth"
	chat "github.com/saufiroja/fin-ai/internal/controllers/chat"
	user "github.com/saufiroja/fin-ai/internal/controllers/user"
	"github.com/saufiroja/fin-ai/internal/middleware"
	"github.com/saufiroja/fin-ai/internal/repositories"
	"github.com/saufiroja/fin-ai/internal/services"
	"github.com/saufiroja/fin-ai/internal/utils"
	"github.com/saufiroja/fin-ai/pkg/databases"
	"github.com/saufiroja/fin-ai/pkg/llm"
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
	llmClient := llm.NewOpenAI(conf)
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
	chatRepository := repositories.NewChatRepository(postgresInstance)
	modelRegistryRepository := repositories.NewModelRegistryRepository(postgresInstance)
	logMessageRepository := repositories.NewLogMessageRepository(postgresInstance)

	authService := services.NewAuthService(userRepository, logger, tokenGenerator, conf)
	userService := services.NewUserService(userRepository, logger)
	logMessageService := services.NewLogMessageService(logMessageRepository, logger)
	chatService := services.NewChatService(
		chatRepository,
		logger,
		llmClient,
		modelRegistryRepository,
		logMessageService,
	)

	authController := authController.NewAuthController(authService, validator)
	userController := user.NewUserController(userService)
	chatController := chat.NewChatController(chatService, validator)

	auth := globalApi.Group("/auth")
	auth.Post("/register", authController.RegisterUser)
	auth.Post("/login", authController.LoginUser)
	auth.Post("/logout", authController.LogoutUser)
	auth.Post("/refresh-token", authController.RefreshToken)

	user := globalApi.Group("/user")
	user.Get("/me", authMiddleware, userController.GetMe)
	user.Put("/:user_id", authMiddleware, userController.UpdateUserById)
	user.Delete("/:user_id", authMiddleware, userController.DeleteUserById)

	chat := globalApi.Group("/chat")
	chat.Post("/session/:user_id", authMiddleware, chatController.CreateChatSession)
	chat.Get("/session/:user_id", authMiddleware, chatController.FindAllChatSessions)
	chat.Put("/session-rename", authMiddleware, chatController.RenameChatSession)
	chat.Delete("/session/:chat_session_id/:user_id", authMiddleware, chatController.DeleteChatSession)
	chat.Post("/send", authMiddleware, chatController.SendChatMessage)
	chat.Get("/session-detail/:chat_session_id/:user_id", authMiddleware, chatController.GetChatSessionDetail)

	if err := a.Listen(fmt.Sprintf("localhost:%s", conf.Http.Port)); err != nil {
		logger.LogPanic(err.Error())
	}
}
