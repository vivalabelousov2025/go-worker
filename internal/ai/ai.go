package ai

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/vivalabelousov2025/go-worker/internal/config"
	"github.com/vivalabelousov2025/go-worker/pkg/logger"
	"google.golang.org/genai"
)

type AiService struct {
	cfg *config.Config
}

func New(cfg *config.Config) *AiService {
	return &AiService{cfg: cfg}
}

func (a *AiService) CallGeminiAPIWithToken(ctx context.Context, prompt string) (string, error) {

	proxyURL, err := url.Parse(a.cfg.ProxyUrl)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, err.Error())
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	clientProxy := &http.Client{
		Transport: transport,
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPClient: clientProxy,
		APIKey:     a.cfg.ApiKey,
		Backend:    genai.BackendGeminiAPI,
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
