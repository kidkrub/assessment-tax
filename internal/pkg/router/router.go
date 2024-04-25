package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func InitRoutes() *echo.Echo {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})
	return e
}
