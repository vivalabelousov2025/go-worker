package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vivalabelousov2025/go-worker/internal/dto"
	"github.com/vivalabelousov2025/go-worker/pkg/logger"
	"go.uber.org/zap"
)

type AiService interface {
	AiRequest(ctx context.Context, prompt string) (string, error)
}

type Handlers struct {
	service AiService
}

func NewHandlers(service AiService) *Handlers {
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

	res, err := h.service.AiRequest(ctx, prompt)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, res)

}

func createPrompt(orders *dto.Order) string {
	prompt := fmt.Sprintf("Составь стэк для разработки по ТЗ: 
						%s, оцени сложность выполнения 
						задания на данном стэке и время выполнение задания в формате:   
						Технологии через запятую   
						Число сложность от 1 до 2   
						Время выполнения число в днях  
						Выведи только нужную информацию и ничего больше",
		orders.Description,
	)

	return prompt
}
