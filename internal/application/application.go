package application

import (
	"github.com/labstack/echo/v4"
	"github.com/vivalabelousov2025/go-worker/internal/heandler"
)

func Run() {
	

	e := echo.New()

	e.POST("/order", heandler.OrderProcessor)

	if err := e.Start(":8082"), err != nil {
		
	}
}