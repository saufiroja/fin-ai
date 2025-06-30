package app

import (
	"fmt"

	"github.com/saufiroja/fin-ai/config"
	"github.com/saufiroja/fin-ai/internal/controllers"
	"github.com/saufiroja/fin-ai/internal/middleware"
	"github.com/saufiroja/fin-ai/internal/repositories"
	"github.com/saufiroja/fin-ai/internal/services"
	"github.com/saufiroja/fin-ai/internal/utils"
	"github.com/saufiroja/fin-ai/pkg/databases"
	"github.com/saufiroja/fin-ai/pkg/llm"
	logging "github.com/saufiroja/fin-ai/pkg/loggings"
	"github.com/saufiroja/fin-ai/pkg/minio"
	"github.com/saufiroja/fin-ai/pkg/redis"
)

type Container struct {
	Dependencies *Dependencies
	Repositories *Repositories
	Services     *Services
	Controllers  *Controllers
}

func NewContainer() *Container {
	container := &Container{}

	container.Dependencies = container.initializeDependencies()
	container.Repositories = container.initializeRepositories()
	container.Services = container.initializeServices()
	container.Controllers = container.initializeControllers()

	return container
}

func (c *Container) initializeDependencies() *Dependencies {
	logger := logging.NewLogrusAdapter()
	conf := config.NewAppConfig(logger)

	postgresInstance := databases.NewPostgres(conf, logger)

	minioClient := minio.NewMinioClient(conf, logger)
	if minioClient == nil {
		logger.LogPanic("failed to create MinIO client")
	}

	redisClient := redis.NewRedisClient(conf, logger)

	openAIClient := llm.NewOpenAI(conf)
	geminiClient := llm.NewGemini(conf)
	validator := utils.NewValidator()
	tokenGenerator := utils.NewJWTTokenGenerator(conf)
	authMiddleware := middleware.Authorization(conf)

	return &Dependencies{
		Logger:         logger,
		Config:         conf,
		Postgres:       postgresInstance,
		Redis:          redisClient,
		MinioClient:    minioClient,
		OpenAIClient:   openAIClient,
		Validator:      validator,
		TokenGen:       tokenGenerator,
		AuthMiddleware: authMiddleware,
		GeminiClient:   geminiClient,
	}
}

func (c *Container) initializeRepositories() *Repositories {
	return &Repositories{
		User:          repositories.NewUserRepository(c.Dependencies.Postgres),
		Chat:          repositories.NewChatRepository(c.Dependencies.Postgres),
		ModelRegistry: repositories.NewModelRegistryRepository(c.Dependencies.Postgres),
		LogMessage:    repositories.NewLogMessageRepository(c.Dependencies.Postgres),
		Transaction:   repositories.NewTransactionRepository(c.Dependencies.Postgres),
		Category:      repositories.NewCategoryRepository(c.Dependencies.Postgres),
		Receipt:       repositories.NewReceiptRepository(c.Dependencies.Postgres),
	}
}

func (c *Container) initializeServices() *Services {
	authService := services.NewAuthService(
		c.Repositories.User,
		c.Dependencies.Logger,
		c.Dependencies.TokenGen,
		c.Dependencies.Config,
	)
	userService := services.NewUserService(c.Repositories.User, c.Dependencies.Logger)
	logMessageService := services.NewLogMessageService(c.Repositories.LogMessage, c.Dependencies.Logger)
	transactionService := services.NewTransactionService(
		c.Repositories.Transaction,
		c.Dependencies.Logger,
		c.Dependencies.OpenAIClient,
	)
	categoryService := services.NewCategoryService(
		c.Repositories.Category,
		c.Dependencies.Logger,
		c.Dependencies.OpenAIClient,
	)
	// Uncomment the following line if you have a Chat service
	chatService := services.NewChatService(
		c.Repositories.Chat,
		c.Dependencies.Logger,
		c.Dependencies.GeminiClient,
		c.Repositories.ModelRegistry,
		c.Repositories.LogMessage,
		transactionService,
		categoryService,
	)
	receiptService := services.NewReceiptService(
		c.Repositories.Receipt,
		transactionService,
		logMessageService,
		categoryService,
		c.Dependencies.MinioClient,
		c.Dependencies.Logger,
		c.Dependencies.OpenAIClient,
		c.Dependencies.GeminiClient,
	)

	return &Services{
		Auth:        authService,
		User:        userService,
		LogMessage:  logMessageService,
		Chat:        chatService,
		Transaction: transactionService,
		Category:    categoryService,
		Receipt:     receiptService,
	}
}

func (c *Container) initializeControllers() *Controllers {
	return &Controllers{
		Auth:        controllers.NewAuthController(c.Services.Auth, c.Dependencies.Validator),
		User:        controllers.NewUserController(c.Services.User),
		Chat:        controllers.NewChatController(c.Services.Chat, c.Dependencies.Validator),
		Transaction: controllers.NewTransactionController(c.Services.Transaction),
		Category:    controllers.NewCategoryController(c.Services.Category),
		Receipt:     controllers.NewReceiptController(c.Services.Receipt),
	}
}

func (c *Container) GetServerAddress() string {
	return fmt.Sprintf("localhost:%s", c.Dependencies.Config.Http.Port)
}
