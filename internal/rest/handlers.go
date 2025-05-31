package rest

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/vivalabelousov2025/go-worker/internal/ai"
	"github.com/vivalabelousov2025/go-worker/internal/calc"
	"github.com/vivalabelousov2025/go-worker/internal/dto"
	"github.com/vivalabelousov2025/go-worker/pkg/logger"
	"go.uber.org/zap"
)

type Handlers struct {
	service *ai.AiService
}

func NewHandlers(service *ai.AiService) *Handlers {
	return &Handlers{service: service}
}

func (h *Handlers) OrderProcess(c echo.Context) error {

	ctx := c.Request().Context()
	var reqSturct dto.Order

	if err := c.Bind(&reqSturct); err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "error", zap.Error(err))
		return err
	}

	var resp dto.Response

	prompt := createPrompt(&reqSturct)

	res, err := h.service.CallGeminiAPIWithToken(ctx, prompt)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "filed to send req Gemeni", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, "filde to send request Gemeni")
	}

	if err := calc.CalcTeam(ctx, reqSturct.Teams, &resp); err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, err.Error())
		return c.JSON(http.StatusBadRequest, "filed to parse date and team")
	}

	return c.JSON(http.StatusOK, res)

}

func createPrompt(orders *dto.Order) string {
	prompt := fmt.Sprintf(`Составь стэк для разработки по ТЗ: Add commentMore actions
						%s, оцени сложность выполнения 
						задания на данном стэке и время выполнение задания в формате:   
						Технологии через запятую   
						Число сложность от 1 до 2   
						Время выполнения число в днях  
						Выведи только нужную информацию и ничего больше`,
		orders.Description,
	)

	return prompt
}

func parseLastTwoNumbers(s string, resp *dto.Response) (int, error) {
	lines := strings.Split(s, "\n")

	const dateFormat = "2006-01-02"

	var numbers []int

	for _, line := range lines {
		// Убираем пробелы по краям, если есть
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue // Пропускаем пустые строки
		}

		// Пытаемся конвертировать строку в число
		num, err := strconv.Atoi(trimmedLine)
		if err == nil {
			// Если успешно, добавляем в наш список чисел
			numbers = append(numbers, num)
		}
		// Если ошибка (не число), просто игнорируем эту строку
	}

	// Проверяем, достаточно ли чисел мы нашли
	if len(numbers) < 2 {
		return 0, fmt.Errorf("не удалось найти два числа, найдено: %d", len(numbers))
	}

	endTime, err := time.Parse(dateFormat, resp.DateStart)

	if err != nil {
		log.Println(err)
		return 0, err
	}

	resp.DateEnd = endTime.AddDate(0, 0, numbers[len(numbers)-1]).GoString()
	// Возвращаем два последних числа
	return numbers[len(numbers)-2], nil
}
