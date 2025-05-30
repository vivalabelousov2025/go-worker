package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handlers struct {
}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (h *Handlers) OrderProcess(c echo.Context) error {

	return c.JSON(http.StatusOK, "test")

}
