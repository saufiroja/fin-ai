package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saufiroja/fin-ai/config"
	"github.com/saufiroja/fin-ai/internal/domains/auth"
	"github.com/saufiroja/fin-ai/internal/domains/categories"
	"github.com/saufiroja/fin-ai/internal/domains/chat"
	"github.com/saufiroja/fin-ai/internal/domains/log_message"
	"github.com/saufiroja/fin-ai/internal/domains/model_registry"
	"github.com/saufiroja/fin-ai/internal/domains/receipt"
	"github.com/saufiroja/fin-ai/internal/domains/transaction"
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
	MinioClient    minio.MinioManager
	OpenAIClient   llm.OpenAI
	Validator      utils.Validator
	TokenGen       utils.TokenGenerator
	AuthMiddleware fiber.Handler
	GeminiClient   llm.Gemini
}

type Repositories struct {
	User          user.UserStorer
	Chat          chat.ChatStorer
	ModelRegistry model_registry.ModelRegistryStorer
	LogMessage    log_message.LogMessageStorer
	Transaction   transaction.TransactionStorer
	Category      categories.CategoryStorer
	Receipt       receipt.ReceiptStorer
}

type Services struct {
	Auth        auth.AuthManager
	User        user.UserManager
	Chat        chat.ChatManager
	LogMessage  log_message.LogMessageManager
	Transaction transaction.TransactionManager
	Category    categories.CategoryManager
	Receipt     receipt.ReceiptManager
}

type Controllers struct {
	Auth        auth.AuthController
	User        user.UserController
	Chat        chat.ChatController
	Transaction transaction.TransactionController
	Category    categories.CategoryController
	Receipt     receipt.ReceiptController
}
