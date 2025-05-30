package rest

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vivalabelousov2025/go-worker/internal/config"
	"github.com/vivalabelousov2025/go-worker/pkg/logger"
	"go.uber.org/zap"
)

type Router struct {
	router   *echo.Echo
	handlers *Handlers
	cfg      *config.Config
}

func NewRouter(ctx context.Context, cfg *config.Config, handlers *Handlers) *Router {

	e := echo.New()

	e.Server.BaseContext = func(_ net.Listener) context.Context {
		return ctx
	}

	e.POST("/order-process", handlers.OrderProcess)

	return &Router{router: e, cfg: cfg, handlers: handlers}

}

func (r *Router) Run(ctx context.Context) {

	restAddr := fmt.Sprintf(":%d", r.cfg.Port)

	if err := r.router.Start(restAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to start server", zap.Error(err))
	}

}
