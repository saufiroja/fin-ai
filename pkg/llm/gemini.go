package llm

import (
	"context"

	"github.com/saufiroja/fin-ai/config"
	"github.com/saufiroja/fin-ai/internal/contracts/responses"
	"github.com/saufiroja/fin-ai/internal/domains/categories"
	"github.com/saufiroja/fin-ai/internal/domains/transaction"
	"github.com/saufiroja/fin-ai/pkg/llm/agents"
	"github.com/saufiroja/fin-ai/pkg/llm/tools"
	"google.golang.org/genai"
)

type Gemini interface {
	Run(ctx context.Context, modelName string, messages []*genai.Content) (*responses.ResponseAI, error)
	RunAgent(ctx context.Context, message string, userId string) (*responses.ResponseAI, error)
	SetTransactionService(transactionService transaction.TransactionManager)
	SetCategoryService(categoryService categories.CategoryManager)
}

type GeminiClient struct {
	conf               *config.AppConfig
	transactionService transaction.TransactionManager
	categoryService    categories.CategoryManager
	transactionAgent   agents.Agent
}

func NewGemini(conf *config.AppConfig) Gemini {
	return &GeminiClient{
		conf:             conf,
		transactionAgent: agents.NewTransactionAgent(conf),
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

func (g *GeminiClient) RunAgent(ctx context.Context, message string, userId string) (*responses.ResponseAI, error) {
	// Create tool context
	toolCtx := &tools.ToolContext{
		TransactionService: g.transactionService,
		CategoryService:    g.categoryService,
		UserId:             userId,
	}

	// Execute using transaction agent
	return g.transactionAgent.Execute(ctx, message, toolCtx)
}

func (g *GeminiClient) SetTransactionService(transactionService transaction.TransactionManager) {
	g.transactionService = transactionService
}

func (g *GeminiClient) SetCategoryService(categoryService categories.CategoryManager) {
	g.categoryService = categoryService
}
