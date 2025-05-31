package rest

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	airequest "github.com/vivalabelousov2025/go-worker/internal/ai"
	"github.com/vivalabelousov2025/go-worker/internal/dto"
	"github.com/vivalabelousov2025/go-worker/pkg/logger"
	"go.uber.org/zap"
)

type Handlers struct {
	service *airequest.AiService
}

func NewHandlers(service *airequest.AiService) *Handlers {
	return &Handlers{service: service}
}

func (h *Handlers) OrderProcess(c echo.Context) error {

	ctx := c.Request().Context()
	var reqSturct dto.Order

	if err := c.Bind(&reqSturct); err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "error", zap.Error(err))
		return err
	}

	prompt := createPrompt(&reqSturct)

	res, err := h.service.AiRequest(prompt)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, res)

}

func createPrompt(orders *dto.Order) string {
	prompt := fmt.Sprintf("собери стэк технлогий для разработки: %s и расчитай коэфицент сложности разработки на данном стэке. Ответ выведи только стэк через запятую и коэфицент в формате числа",
		orders.Description,
	)

	return prompt
}
