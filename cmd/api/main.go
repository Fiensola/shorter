package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"shorter/internal/app"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	app := app.NewApp()
	defer app.Db.Close()

	// start server
	go func() {
		if err := app.StartHttpServer(); err != nil && err != http.ErrServerClosed {
			app.Logger.Fatal("server startup failed", zap.Error(err))
		}
	}()

	// wait stop signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	app.Logger.Info("shutting down server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownHttpServer(ctx); err != nil {
		app.Logger.Fatal("server shutdown failed", zap.Error(err))
	}

	app.Logger.Info("server stopped gracefully")
}
