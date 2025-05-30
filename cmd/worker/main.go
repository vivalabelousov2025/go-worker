package main

import (
	"context"

	"github.com/vivalabelousov2025/go-worker/internal/application"
	"github.com/vivalabelousov2025/go-worker/internal/config"
	"github.com/vivalabelousov2025/go-worker/pkg/logger"
)

func main() {

	ctx, _ := logger.New(context.Background())

	cfg, err := config.New()

	if err != nil {
		return
	}

	application.Run(ctx, cfg)
}
