package main

import (
	"fmt"
	"net/http"

	"github.com/kidkrub/assessment-tax/internal/pkg/config"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg := config.New()
	serverConfig := cfg.Server()
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})
	server := fmt.Sprintf("%s:%d", serverConfig.Hostname, serverConfig.PORT)
	e.Logger.Fatal(e.Start(server))
}
