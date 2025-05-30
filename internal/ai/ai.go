package airequest

import (
	"context"

	"github.com/vivalabelousov2025/go-worker/internal/config"
	"google.golang.org/genai"
)

type AiService struct {
}

func New() *AiService {
	return &AiService{}
}

func (a *AiService) AiRequest(prompt string, cfg *config.Config) (string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  cfg.ApiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", err
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", err
	}
	res := result.Text()
	return res, nil
}
