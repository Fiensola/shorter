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
	defer app.Logger.Sync()

	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel() // if panic

	// start server
	go func() {
		if err := app.StartHttpServer(); err != nil && err != http.ErrServerClosed {
			app.Logger.Error("server startup failed", zap.Error(err))
			cancel()
		}
	}()

	// start consumer
	go func() {
		if err := app.KafkaConsumer.Start(rootCtx); err != nil {
			app.Logger.Error("kafka consumer error", zap.Error(err))
			cancel()
		}
	}()

	// wait stop signal
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	app.Logger.Info("shutting down server ...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := app.ShutdownHttpServer(shutdownCtx); err != nil {
		app.Logger.Error("server shutdown failed", zap.Error(err))
	}

	if err := app.KafkaProducer.Close(); err != nil {
		app.Logger.Error("Kafka producer close error", zap.Error(err))
	}

	app.Logger.Info("server stopped gracefully")
}
