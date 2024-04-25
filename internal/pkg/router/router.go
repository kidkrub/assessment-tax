package router

import (
	"net/http"

	"github.com/kidkrub/assessment-tax/internal/pkg/handler/tax"
	"github.com/labstack/echo/v4"
)

func InitRoutes() *echo.Echo {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})
	th := tax.New()

	e.POST("/tax/calculations", th.TaxCalculateHandler)

	return e
}
