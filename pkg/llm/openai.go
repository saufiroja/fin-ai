package llm

import (
	"context"
	"sync"

	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/saufiroja/fin-ai/config"
	"github.com/saufiroja/fin-ai/internal/models"
)

type OpenAI interface {
	SendChat(ctx context.Context, modelName string, messages []openai.ChatCompletionMessageParamUnion) (*models.ResponseAI, error)
}

type OpenAIClient struct {
	client openai.Client
	model  string
}

var (
	instance *OpenAIClient
	once     sync.Once
)

func NewOpenAI(conf *config.AppConfig) OpenAI {
	once.Do(func() {
		instance = &OpenAIClient{
			client: openai.NewClient(
				option.WithAPIKey(conf.OpenAI.ApiKey),
			),
		}
	})
	return instance
}

func (o *OpenAIClient) SendChat(ctx context.Context, modelName string, messages []openai.ChatCompletionMessageParamUnion) (*models.ResponseAI, error) {
	resp, err := o.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    modelName,
		Messages: messages,
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, nil
	}

	res := &models.ResponseAI{
		Response:    resp.Choices[0].Message.Content,
		InputToken:  int(resp.Usage.PromptTokens),
		OutputToken: int(resp.Usage.CompletionTokens),
	}

	return res, nil
}
