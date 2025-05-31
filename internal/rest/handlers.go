package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/vivalabelousov2025/go-worker/internal/ai"
	"github.com/vivalabelousov2025/go-worker/internal/calc"
	"github.com/vivalabelousov2025/go-worker/internal/config"
	"github.com/vivalabelousov2025/go-worker/internal/dto"
	"github.com/vivalabelousov2025/go-worker/pkg/logger"
	"go.uber.org/zap"
)

type Handlers struct {
	service *ai.AiService
	cfg     *config.Config
}

func NewHandlers(service *ai.AiService, cfg *config.Config) *Handlers {
	return &Handlers{service: service, cfg: cfg}
}

func (h *Handlers) OrderProcess(c echo.Context) error {
	ctx := c.Request().Context()
	var reqSturct dto.Order

	// Логируем тело запроса
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "ошибка чтения тела запроса", zap.Error(err))
		return err
	}
	fmt.Printf("Получено тело запроса: %s\n", string(body))

	// Восстанавливаем тело запроса для повторного чтения
	c.Request().Body = io.NopCloser(bytes.NewBuffer(body))

	if err := c.Bind(&reqSturct); err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "ошибка биндинга", zap.Error(err))
		return err
	}

	fmt.Printf("Получен запрос: %+v\n", reqSturct)
	fmt.Printf("OrderID: %s\n", reqSturct.OrderID)
	fmt.Printf("UserID: %s\n", reqSturct.UserID)
	fmt.Printf("Description: %s\n", reqSturct.Description)

	var resp dto.Response

	for i := 1; i < 5; i++ {
		team, err := h.GetTeam(c, &resp)
		technologies, err := h.GetTechnologies()
		prompt := createPrompt(&reqSturct, &technologies)
		fmt.Println(team, 44)

		res, err := h.service.CallGeminiAPIWithToken(ctx, prompt)
		if err != nil {
			logger.GetLoggerFromCtx(ctx).Info(ctx, "filed to send req Gemeni", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, "filde to send request Gemeni")
		}

		tech, hard, err := parseLastTwoNumbers(res, &resp)
		if err != nil {
			logger.GetLoggerFromCtx(ctx).Info(ctx, "ошибка при парсинге ответа", zap.Error(err))
			continue
		}

		if team == nil {
			logger.GetLoggerFromCtx(ctx).Info(ctx, "команда не найдена")
			continue
		}

		resp.Price, err = calc.CalcPrice(&reqSturct, team, hard)
		if err != nil {
			logger.GetLoggerFromCtx(ctx).Info(ctx, "ошибка при расчете цены", zap.Error(err))
			continue
		}

		resp.OrderID = reqSturct.OrderID
		resp.Stack = tech

		err = resp.Validate()

		if err != nil {
			continue
		}
		break
	}

	err = h.UpdateOrder(c, &reqSturct, &resp)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "ошибка при обновлении заказа", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, "ошибка при обновлении заказа")
	}
	return c.JSON(http.StatusOK, resp)
}

func createPrompt(orders *dto.Order, technologies *[]dto.OrderTechnology) string {
	prompt := fmt.Sprintf(`Составь стэк для разработки по ТЗ:

						%s, оцени сложность выполнения

						 задания на данном стэке и время выполнение задания в формате:
						 Без текста в начале


						 Число сложность от 1 до 2 без текста


						 Время выполнения число в днях  без текста

					%s Выбери из этого списка айди нужных технологий и верни их uuid через запятую без лишнего 
					`,
		orders.Description,
		technologies,
	)

	return prompt
}

func parseLastTwoNumbers(s string, resp *dto.Response) ([]string, float64, error) {

	const dateFormat = "2006-01-02"

	str := strings.TrimSpace(s)

	str = strings.ReplaceAll(str, "\n", ",")

	arr := strings.Split(str, ",")
	fmt.Println(resp.DateStart, "resp date start")

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

func (h *Handlers) GetTechnologies() ([]dto.OrderTechnology, error) {
	resp, err := http.Get(fmt.Sprintf("%s/technologies", h.cfg.BackendUrl))
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе технологий: %v", err)
	}
	defer resp.Body.Close()

	var technologies []dto.OrderTechnology
	if err := json.NewDecoder(resp.Body).Decode(&technologies); err != nil {
		return nil, fmt.Errorf("ошибка при обработке данных технологий: %v", err)
	}

	return technologies, nil
}

func (h *Handlers) GetTeam(c echo.Context, response *dto.Response) (*dto.Team, error) {
	resp, err := http.Get(fmt.Sprintf("%s/teams-by-worker", h.cfg.BackendUrl))
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе: %v", err)
	}
	defer resp.Body.Close()

	var teams []dto.Team
	if err := json.NewDecoder(resp.Body).Decode(&teams); err != nil {
		return nil, fmt.Errorf("ошибка при обработке данных команд: %v", err)
	}

	if teams == nil || len(teams) == 0 {
		return nil, fmt.Errorf("нет доступных команд")
	}

	fmt.Println(teams)
	team, err := calc.CalcTeam(c.Request().Context(), teams, response)
	if err != nil {
		return nil, fmt.Errorf("ошибка при обработке данных команды: %v", err)
	}

	return team, nil
}

func (h *Handlers) UpdateOrder(c echo.Context, order *dto.Order, resp *dto.Response) error {
	fmt.Printf("Отправляем запрос на обновление заказа: %+v\n", order)

	// Добавляем обязательные поля из resp
	order.TeamID = resp.TeamID
	order.TotalPrice = resp.Price

	// Фильтруем пустые строки из Stack
	var filteredStack []string
	for _, tech := range resp.Stack {
		if strings.TrimSpace(tech) != "" {
			filteredStack = append(filteredStack, tech)
		}
	}

	// Создаем структуру для запроса с order_technologies
	requestData := struct {
		*dto.Order
		OrderTechnologies []string `json:"order_technologies"`
	}{
		Order:             order,
		OrderTechnologies: filteredStack,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("ошибка при сериализации: %v", err)
	}
	fmt.Printf("Подготовленные JSON данные: %s\n", string(jsonData))

	url := fmt.Sprintf("%s/worker/order-process", h.cfg.BackendUrl)
	fmt.Printf("Отправляем POST запрос на %s\n", url)

	httpResp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Ошибка при отправке запроса: %v\n", err)
		return fmt.Errorf("ошибка при запросе: %v", err)
	}
	defer httpResp.Body.Close()

	fmt.Printf("Получен ответ со статусом: %d\n", httpResp.StatusCode)

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("ошибка при чтении ответа: %v", err)
	}
	fmt.Printf("Тело ответа: %s\n", string(body))

	return nil
}
