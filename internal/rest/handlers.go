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

	for i := 1; i < 5; i++ {

		res, err := h.service.CallGeminiAPIWithToken(ctx, prompt)
		if err != nil {
			logger.GetLoggerFromCtx(ctx).Info(ctx, "filed to send req Gemeni", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, "filde to send request Gemeni")
		}

		team, err := calc.CalcTeam(ctx, reqSturct.Teams, &resp)
		if err != nil {
			logger.GetLoggerFromCtx(ctx).Info(ctx, err.Error())
			return c.JSON(http.StatusBadRequest, "filed to parse date and team")
		}

		tech, hard, err := parseLastTwoNumbers(res, &resp)

		resp.Price, err = calc.CalcPrice(&reqSturct, team, hard)

		resp.OrderID = reqSturct.OrderID
		resp.Stack = tech

		err = resp.Validate()

		if err != nil {
			continue
		}
		break
	}

	return c.JSON(http.StatusOK, resp)

}

func createPrompt(orders *dto.Order) string {
	prompt := fmt.Sprintf(`Составь стэк для разработки по ТЗ:

						%s, оцени сложность выполнения

						 задания на данном стэке и время выполнение задания в формате:
						 Без текста в начале

						 Технологии в одну строку через запятую без лишнего (не более 7)

						 Число сложность от 1 до 2 без текста

						 Время выполнения число в днях  без текста

						Выведи только технологии и числа

						 пример:

						REACT, JS, GO

						1.6

						 20`,
		orders.Description,
	)

	return prompt
}

func parseLastTwoNumbers(s string, resp *dto.Response) ([]string, float64, error) {

	const dateFormat = "2006-01-02"

	str := strings.TrimSpace(s)

	str = strings.ReplaceAll(str, "\n", ",")

	arr := strings.Split(str, ",")
	fmt.Println(arr)

	endTime, err := time.Parse(dateFormat, resp.DateStart)

	if err != nil {
		log.Println(err)
		return nil, 0, err
	}

	endDate, _ := strconv.Atoi(arr[len(arr)-1])

	resp.DateEnd = endTime.AddDate(0, 0, endDate).Format(dateFormat)

	hard, err := strconv.ParseFloat(arr[len(arr)-1], 32)
	if err != nil {
		log.Print(err)
	}

	return arr[:len(arr)-3], hard, nil
}
