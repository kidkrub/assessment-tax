package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kidkrub/assessment-tax/internal/pkg/config"
	"github.com/kidkrub/assessment-tax/internal/pkg/router"
)

func main() {
	cfg := config.New()
	serverConfig := cfg.Server()

	server := fmt.Sprintf("%s:%d", serverConfig.Hostname, serverConfig.PORT)

	e := router.InitRoutes()

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
