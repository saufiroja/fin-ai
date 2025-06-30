package llm

import (
	"context"

	"github.com/saufiroja/fin-ai/config"
	"github.com/saufiroja/fin-ai/internal/contracts/responses"
	"google.golang.org/genai"
)

type Gemini interface {
	Run(ctx context.Context, modelName string, messages []*genai.Content) (*responses.ResponseAI, error)
}

type GeminiClient struct {
	conf *config.AppConfig
}

func NewGemini(conf *config.AppConfig) Gemini {
	return &GeminiClient{
		conf: conf,
	}
}

func (g *GeminiClient) Run(ctx context.Context, modelName string, messages []*genai.Content) (*responses.ResponseAI, error) {
	client, err := genai.NewClient(
		ctx,
		&genai.ClientConfig{
			APIKey: g.conf.Gemini.ApiKey,
		},
	)
	if err != nil {
		return nil, err
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		messages,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &responses.ResponseAI{
		Response:    result.Text(),
		InputToken:  int(result.UsageMetadata.PromptTokenCount),
		OutputToken: int(result.UsageMetadata.CandidatesTokenCount),
	}, nil
}
