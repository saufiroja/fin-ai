package llm

import (
	"context"
	"fmt"
	"strings"
	"sync"

	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/ssestream"
	"github.com/saufiroja/fin-ai/config"
	"github.com/saufiroja/fin-ai/internal/contracts/responses"
)

type OpenAI interface {
	SendChat(ctx context.Context, modelName string, messages []openai.ChatCompletionMessageParamUnion) (*responses.ResponseAI, error)
	SendChatStream(ctx context.Context, modelName string, messages string) (*ssestream.Stream[openai.ChatCompletionChunk], error)
	CreateEmbedding(ctx context.Context, input openai.EmbeddingNewParamsInputUnion) *responses.ResponseEmbedding
}

type OpenAIClient struct {
	client openai.Client
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

func (o *OpenAIClient) SendChat(ctx context.Context, modelName string, messages []openai.ChatCompletionMessageParamUnion) (*responses.ResponseAI, error) {
	// Enhanced parameters for better accuracy
	var model openai.ChatModel
	switch modelName {
	case "gpt-4o":
		model = openai.ChatModelGPT4o
	case "gpt-4o-mini":
		model = openai.ChatModelGPT4oMini
	default:
		model = openai.ChatModelGPT4o // Default to GPT-4o for better accuracy
	}

	resp, err := o.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:       model,
		Messages:    messages,
		Temperature: openai.Float(0.0), // Zero temperature for maximum consistency
		MaxTokens:   openai.Int(4000),  // Optimized for receipt data
		Seed:        openai.Int(12345), // Fixed seed for reproducible results
		TopP:        openai.Float(0.1), // Low top_p for more focused responses
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, nil
	}

	res := &responses.ResponseAI{
		Response:    resp.Choices[0].Message.Content,
		InputToken:  int(resp.Usage.PromptTokens),
		OutputToken: int(resp.Usage.CompletionTokens),
	}

	return res, nil
}

func (o *OpenAIClient) SendChatStream(ctx context.Context, modelName string, messages string) (*ssestream.Stream[openai.ChatCompletionChunk], error) {
	stream := o.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(messages),
		},
		Seed:  openai.Int(0),
		Model: openai.ChatModelGPT4o,
	})

	return stream, nil
}

func (o *OpenAIClient) CreateEmbedding(ctx context.Context, input openai.EmbeddingNewParamsInputUnion) *responses.ResponseEmbedding {
	resp, err := o.client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Model: "text-embedding-3-small",
		Input: input,
	})
	if err != nil {
		return nil
	}

	if len(resp.Data) == 0 {
		return nil
	}

	// float64 to pgvector
	embedding := make([]string, len(resp.Data[0].Embedding))
	for i, v := range resp.Data[0].Embedding {
		embedding[i] = fmt.Sprintf("%f", v)
	}

	embeddingData := fmt.Sprintf("[%s]", strings.Join(embedding, ","))
	res := &responses.ResponseEmbedding{
		Embeddings:  embeddingData,
		InputToken:  int(resp.Usage.PromptTokens),
		OutputToken: int(resp.Usage.TotalTokens),
	}

	return res
}
