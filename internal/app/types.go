package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saufiroja/fin-ai/config"
	"github.com/saufiroja/fin-ai/internal/domains/auth"
	"github.com/saufiroja/fin-ai/internal/domains/chat"
	"github.com/saufiroja/fin-ai/internal/domains/log_message"
	"github.com/saufiroja/fin-ai/internal/domains/model_registry"
	"github.com/saufiroja/fin-ai/internal/domains/user"
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
	User          user.UserRepository
	Chat          chat.ChatRepository
	ModelRegistry model_registry.ModelRegistryRepository
	LogMessage    log_message.LogMessageRepository
}

type Services struct {
	Auth       auth.AuthService
	User       user.UserService
	Chat       chat.ChatService
	LogMessage log_message.LogMessageService
}

type Controllers struct {
	Auth auth.AuthController
	User user.UserController
	Chat chat.ChatController
}
