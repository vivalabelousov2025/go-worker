package heandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func OrderProcessor(c echo.Context) error {

	return c.JSON(http.StatusOK, "....")

}
