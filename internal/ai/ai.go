package airequest

import (
	"context"
	"fmt"
	"log"

	"github.com/vivalabelousov2025/go-worker/internal/config"
	"google.golang.org/genai"
)

type AiService struct {
	cfg *config.Config
}

func New(cfg *config.Config) *AiService {
	return &AiService{cfg: cfg}
}

func (a *AiService) AiRequest(prompt string) (string, error) {
	fmt.Println("ai request")
	ctx := context.Background()
	fmt.Println(a.cfg.ApiKey)
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  a.cfg.ApiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Text())
	return result.Text(), nil
}
