package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	ctx, shutdown := signal.NotifyContext(context.Background(), os.Interrupt)
	defer shutdown()

	go func() {
		if err := e.Start(server); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}
