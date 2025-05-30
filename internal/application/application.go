package application

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vivalabelousov2025/go-worker/internal/config"
	"github.com/vivalabelousov2025/go-worker/internal/heandler"
	"github.com/vivalabelousov2025/go-worker/pkg/logger"
)

func Run(ctx context.Context, cfg *config.Config) {

	e := echo.New()

	e.POST("/order", heandler.OrderProcessor)

	logger.GetLoggerFromCtx(ctx).Info(ctx, "starting server")

	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.GetLoggerFromCtx(ctx).Info(ctx, err.Error())
	}
}
