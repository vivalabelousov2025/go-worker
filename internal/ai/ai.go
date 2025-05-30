package airequest

import (
	"context"
	"fmt"

	"github.com/vivalabelousov2025/go-worker/internal/config"
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
		return "", err
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		genai.Text(prompt),
		nil,
	)

	fmt.Println(result.Text(), err)
	if err != nil {
		return "", err
	}
	res := result.Text()
	fmt.Println(res)
	return res, nil
}
