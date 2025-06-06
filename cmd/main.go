package main

import (
	"context"

	"github.com/vivalabelousov2025/go-worker/internal/ai"
	"github.com/vivalabelousov2025/go-worker/internal/config"
	"github.com/vivalabelousov2025/go-worker/internal/rest"
	"github.com/vivalabelousov2025/go-worker/pkg/logger"
)

func main() {

	ctx, _ := logger.New(context.Background())

	cfg, err := config.New()
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "filed to parse config")
	}

	service := ai.New(cfg)

	handl := rest.NewHandlers(service, cfg)

	router := rest.NewRouter(ctx, cfg, handl)

	router.Run(ctx)

}
