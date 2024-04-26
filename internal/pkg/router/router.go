package router

import (
	"database/sql"
	"net/http"

	"github.com/kidkrub/assessment-tax/internal/pkg/handler/admin"
	"github.com/kidkrub/assessment-tax/internal/pkg/handler/tax"
	"github.com/labstack/echo/v4"
)

func InitRoutes(db *sql.DB) *echo.Echo {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})
	th := tax.New()
	ah := admin.New(db)

	e.POST("/tax/calculations", th.TaxCalculateHandler)
	e.POST("/admin/deductions/:type", ah.SetDeductionValueHandler)

	return e
}
