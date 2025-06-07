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

	llmClient := llm.NewOpenAI(conf)
	validator := utils.NewValidator()
	tokenGenerator := utils.NewJWTTokenGenerator(conf)
	authMiddleware := middleware.Authorization(conf)

	return &Dependencies{
		Logger:         logger,
		Config:         conf,
		Postgres:       postgresInstance,
		Redis:          redisClient,
		MinioClient:    minioClient,
		LLMClient:      llmClient,
		Validator:      validator,
		TokenGen:       tokenGenerator,
		AuthMiddleware: authMiddleware,
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
	}
}

func (c *Container) initializeServices() *Services {
	logMessageService := services.NewLogMessageService(c.Repositories.LogMessage, c.Dependencies.Logger)

	return &Services{
		Auth: services.NewAuthService(
			c.Repositories.User,
			c.Dependencies.Logger,
			c.Dependencies.TokenGen,
			c.Dependencies.Config,
		),
		User:       services.NewUserService(c.Repositories.User, c.Dependencies.Logger),
		LogMessage: logMessageService,
		// Chat: services.NewChatService(
		// 	c.Repositories.Chat,
		// 	c.Dependencies.Logger,
		// 	c.Dependencies.LLMClient,
		// 	c.Repositories.ModelRegistry,
		// 	logMessageService,
		// ),
		Transaction: services.NewTransactionService(
			c.Repositories.Transaction,
			c.Dependencies.Logger,
			c.Dependencies.LLMClient,
		),
		Category: services.NewCategoryService(
			c.Repositories.Category,
			c.Dependencies.Logger,
			c.Dependencies.LLMClient,
		),
	}
}

func (c *Container) initializeControllers() *Controllers {
	return &Controllers{
		Auth:        controllers.NewAuthController(c.Services.Auth, c.Dependencies.Validator),
		User:        controllers.NewUserController(c.Services.User),
		Chat:        controllers.NewChatController(c.Services.Chat, c.Dependencies.Validator),
		Transaction: controllers.NewTransactionController(c.Services.Transaction),
		Category:    controllers.NewCategoryController(c.Services.Category),
	}
}

func (c *Container) GetServerAddress() string {
	return fmt.Sprintf("localhost:%s", c.Dependencies.Config.Http.Port)
}
