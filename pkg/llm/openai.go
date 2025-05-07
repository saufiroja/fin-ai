package llm

import (
	"context"
	"sync"

	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/saufiroja/fin-ai/config"
)

type OpenAIClient struct {
	client openai.Client
	model  string
}

var (
	instance *OpenAIClient
	once     sync.Once
)

func NewOpenAI(conf *config.AppConfig) *OpenAIClient {
	once.Do(func() {
		instance = &OpenAIClient{
			client: openai.NewClient(
				option.WithAPIKey(conf.OpenAI.ApiKey),
			),
			model: conf.OpenAI.Model,
		}
	})
	return instance
}

func (o *OpenAIClient) Chat(ctx context.Context, messages []openai.ChatCompletionMessageParamUnion) (string, error) {
	resp, err := o.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    o.model,
		Messages: messages,
	})
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", nil
	}

	return resp.Choices[0].Message.Content, nil
}
