package airequest

import (
	"context"
	"fmt"

	"github.com/vivalabelousov2025/go-worker/internal/config"
	"github.com/vivalabelousov2025/go-worker/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/genai"
)

type AiService struct {
	cfg *config.Config
}

func New(cfg *config.Config) *AiService {
	return &AiService{cfg: cfg}
}

func (a *AiService) AiRequest(ctx context.Context, prompt string) (string, error) {
	fmt.Println("ai request")

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  a.cfg.ApiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "filed connect to gemeni", zap.Error(err))
		return "", err
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		genai.Text(prompt),
		nil,
	)

	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "filed to generate response", zap.Error(err))
		return "", err
	}
	res := result.Text()

	return res, nil
}
