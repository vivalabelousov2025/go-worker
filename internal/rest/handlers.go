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

	_, _, err = parseLastTwoNumbers(res, &resp)

	return c.JSON(http.StatusOK, resp)

}

func createPrompt(orders *dto.Order) string {
	prompt := fmt.Sprintf(`Составь стэк для разработки по ТЗ: Add commentMore actions
							%s, оцени сложность выполнения 
							задания на данном стэке и время выполнение задания в формате:   
							Технологии через запятую   
							Число сложность от 1 до 2   
							Время выполнения число в днях  
							Выведи только технологии и числа
							пример:
							REACT, JS, GO
							1.6
							20`,
		orders.Description,
	)

	return prompt
}

func parseLastTwoNumbers(s string, resp *dto.Response) (int, []string, error) {
	lines := strings.Split(s, "\n")

	const dateFormat = "2006-01-02"

	stack := strings.Split(lines[0], ",")
	for i, word := range stack {
		stack[i] = strings.TrimSpace(word) // Убираем лишние пробелы
	}

	num1, err1 := strconv.Atoi(strings.TrimSpace(lines[1]))
	num2, err2 := strconv.Atoi(strings.TrimSpace(lines[2]))
	if err1 != nil || err2 != nil {
		fmt.Println("Ошибка парсинга чисел")
		return 0, nil, err1
	}

	fmt.Println(num1, num2, stack)

	endTime, err := time.Parse(dateFormat, resp.DateStart)

	if err != nil {
		log.Println(err)
		return 0, nil, err
	}

	resp.DateEnd = endTime.AddDate(0, 0, num1).Format(dateFormat)
	// Возвращаем два последних числа
	return num2, nil, nil
}
