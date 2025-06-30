package config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
	logging "github.com/saufiroja/fin-ai/pkg/loggings"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type AppConfig struct {
	App struct {
		Env string
	}
	Http struct {
		Port string
	}
	Postgres struct {
		Name string
		User string
		Pass string
		Host string
		Port string
		SSL  string
	}
	Jwt struct {
		Secret string
	}
	OpenAI struct {
		ApiKey string
	}
	Redis struct {
		URL string
	}
	Minio struct {
		Endpoint  string
		AccessKey string
		SecretKey string
		Bucket    string
		Region    string
		UseSSL    bool
	}
	Gemini struct {
		ApiKey string
	}
}

var appConfig *AppConfig
var lock = &sync.Mutex{}

func NewAppConfig(logging logging.Logger) *AppConfig {
	// add config file path in .env
	_ = godotenv.Load("../.env")

	if appConfig == nil {
		lock.Lock()
		defer lock.Unlock()
		if appConfig == nil {
			logging.LogInfo("Creating AppConfig first time")
			appConfig = &AppConfig{}

			appConfig.initApp()
			appConfig.initHttp()
			appConfig.initPostgres()
			appConfig.initJwt()
			appConfig.initOpenAI(logging)
			appConfig.initRedis()
			appConfig.initMinio()
			appConfig.initGemini()
		} else {
			logging.LogInfo("AppConfig already created")
		}
	} else {
		logging.LogInfo("AppConfig already created")
	}

	return appConfig
}

func (c *AppConfig) initApp() {
	c.App.Env = os.Getenv("GO_ENV")
	switch cases.Lower(language.English).String(c.App.Env) {
	case "development":
		c.App.Env = "development"
	case "staging":
		c.App.Env = "staging"
	case "testing":
		c.App.Env = "testing"
	case "production":
		c.App.Env = "production"
	default:
		c.App.Env = "development"
	}
}

func (c *AppConfig) initHttp() {
	c.Http.Port = os.Getenv("HTTP_PORT")
	if c.Http.Port == "" {
		c.Http.Port = "8080"
	}
}

func (c *AppConfig) initPostgres() {
	c.Postgres.Host = os.Getenv("DB_HOST")
	c.Postgres.Port = os.Getenv("DB_PORT")
	c.Postgres.User = os.Getenv("DB_USER")
	c.Postgres.Pass = os.Getenv("DB_PASS")
	c.Postgres.Name = os.Getenv("DB_NAME")
	c.Postgres.SSL = os.Getenv("DB_SSL_MODE")
}

func (c *AppConfig) initJwt() {
	c.Jwt.Secret = os.Getenv("JWT_SECRET")
	if c.Jwt.Secret == "" {
		c.Jwt.Secret = "secret"
	}
}

func (c *AppConfig) initOpenAI(logging logging.Logger) {
	c.OpenAI.ApiKey = os.Getenv("OPENAI_API_KEY")
	if c.OpenAI.ApiKey == "" {
		logging.LogPanic("OpenAI API key not found")
	}
}

func (c *AppConfig) initRedis() {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}

	c.Redis.URL = host + ":" + port
}

func (c *AppConfig) initMinio() {
	c.Minio.Endpoint = os.Getenv("MINIO_ENDPOINT")
	c.Minio.AccessKey = os.Getenv("MINIO_ACCESS_KEY")
	c.Minio.SecretKey = os.Getenv("MINIO_SECRET_KEY")
	c.Minio.Region = os.Getenv("MINIO_REGION")
	c.Minio.UseSSL = os.Getenv("MINIO_USE_SSL") == "true"
}

func (c *AppConfig) initGemini() {
	c.Gemini.ApiKey = os.Getenv("GEMINI_API_KEY")
	if c.Gemini.ApiKey == "" {
		panic("Gemini API key not found")
	}
}
