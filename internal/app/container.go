package app

import (
	"fmt"

	"github.com/saufiroja/fin-ai/config"
	"github.com/saufiroja/fin-ai/internal/controllers/auth"
	"github.com/saufiroja/fin-ai/internal/controllers/chat"
	"github.com/saufiroja/fin-ai/internal/controllers/user"
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
	c.setupPostgresCleanup(postgresInstance, logger)

	minioClient := minio.NewMinioClient(conf, logger)
	if minioClient == nil {
		logger.LogPanic("failed to create MinIO client")
	}

	redisClient := redis.NewRedisClient(conf, logger)
	c.setupRedisCleanup(redisClient, logger)

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
		Chat: services.NewChatService(
			c.Repositories.Chat,
			c.Dependencies.Logger,
			c.Dependencies.LLMClient,
			c.Repositories.ModelRegistry,
			logMessageService,
		),
	}
}

func (c *Container) initializeControllers() *Controllers {
	return &Controllers{
		Auth: auth.NewAuthController(c.Services.Auth, c.Dependencies.Validator),
		User: user.NewUserController(c.Services.User),
		Chat: chat.NewChatController(c.Services.Chat, c.Dependencies.Validator),
	}
}

func (c *Container) GetServerAddress() string {
	return fmt.Sprintf("localhost:%s", c.Dependencies.Config.Http.Port)
}

func (c *Container) setupPostgresCleanup(postgres databases.PostgresManager, logger logging.Logger) {
	defer func() {
		if err := postgres.CloseConnection(); err != nil {
			logger.LogError(fmt.Sprintf("failed to close postgres connection: %v", err))
		}
	}()
}

func (c *Container) setupRedisCleanup(redis *redis.RedisClient, logger logging.Logger) {
	defer func() {
		if err := redis.Close(); err != nil {
			logger.LogError(fmt.Sprintf("failed to close redis connection: %v", err))
		}
	}()
}
