package example

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetData(c echo.Context) error {
	return c.String(http.StatusOK, "OK 111s")
}
