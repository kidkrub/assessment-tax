package router

import (
	"database/sql"
	"net/http"

	"github.com/kidkrub/assessment-tax/internal/pkg/handler/admin"
	"github.com/kidkrub/assessment-tax/internal/pkg/handler/tax"
	cmw "github.com/kidkrub/assessment-tax/internal/pkg/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRoutes(db *sql.DB) *echo.Echo {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})
	th := tax.New()
	ah := admin.New(db)

	ag := e.Group("/admin")
	ag.Use(middleware.BasicAuth(cmw.BasicAuthenticate()))

	e.POST("/tax/calculations", th.TaxCalculateHandler)
	ag.POST("/deductions/:type", ah.SetDeductionValueHandler)

	return e
}
