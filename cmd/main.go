package main

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/openai/openai-go"
	"github.com/saufiroja/fin-ai/config"
	"github.com/saufiroja/fin-ai/pkg/databases"
	"github.com/saufiroja/fin-ai/pkg/llm"
	logging "github.com/saufiroja/fin-ai/pkg/loggings"
)

func main() {
	app := fiber.New()
	logger := logging.NewLogrusAdapter()
	conf := config.NewAppConfig(logger)
	client := llm.NewOpenAI(conf)
	databases.NewPostgres(conf, logger)

	app.Get("/", func(c *fiber.Ctx) error {
		systemPrompt := "You are a helpful assistant."
		messages := []openai.ChatCompletionMessageParamUnion{openai.SystemMessage(systemPrompt)}

		messages = append(messages, openai.UserMessage("What is the capital of France?"))
		response, err := client.Chat(context.Background(), messages)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error: %v", err))
		}

		return c.SendString(fmt.Sprintf("Response: %s", response))
	})

	app.Listen(":8080")
}
