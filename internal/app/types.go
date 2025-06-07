package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saufiroja/fin-ai/config"
	"github.com/saufiroja/fin-ai/internal/controllers/auth"
	"github.com/saufiroja/fin-ai/internal/controllers/chat"
	"github.com/saufiroja/fin-ai/internal/controllers/user"
	"github.com/saufiroja/fin-ai/internal/domains"
	"github.com/saufiroja/fin-ai/internal/utils"
	"github.com/saufiroja/fin-ai/pkg/databases"
	"github.com/saufiroja/fin-ai/pkg/llm"
	logging "github.com/saufiroja/fin-ai/pkg/loggings"
	"github.com/saufiroja/fin-ai/pkg/minio"
	"github.com/saufiroja/fin-ai/pkg/redis"
)

type Dependencies struct {
	Logger         logging.Logger
	Config         *config.AppConfig
	Postgres       databases.PostgresManager
	Redis          *redis.RedisClient
	MinioClient    *minio.MinioClient
	LLMClient      llm.OpenAI
	Validator      utils.Validator
	TokenGen       utils.TokenGenerator
	AuthMiddleware fiber.Handler
}

type Repositories struct {
	User          domains.UserRepositoryInterface
	Chat          domains.ChatRepositoryInterface
	ModelRegistry domains.ModelRegistryRepositoryInterface
	LogMessage    domains.LogMessageRepositoryInterface
}

type Services struct {
	Auth       domains.AuthServiceInterface
	User       domains.UserServiceInterface
	Chat       domains.ChatServiceInterface
	LogMessage domains.LogMessageServiceInterface
}

type Controllers struct {
	Auth auth.AuthController
	User user.UserController
	Chat chat.ChatController
}
