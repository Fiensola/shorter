package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"shorter/internal/config"
	"shorter/internal/logger"
	"shorter/internal/server"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	// load config
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	// init logger
	logger, err := logger.NewLogger(cfg.IsDev)
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}

	// init db
	db, err := connectDb(cfg, logger)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// init server
	srv := server.NewServer(
		cfg.Server.Host+":"+fmt.Sprint(cfg.Server.Port),
		logger,
		cfg,
		db,
	)

	// start server
	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server startup failed", zap.Error(err))
		}
	}()

	// wait stop signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	logger.Info("shutting down server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("server shutdown failed", zap.Error(err))
	}

	logger.Info("server stopped gracefully")
}

func connectDb(c *config.Config, logger *zap.Logger) (*pgxpool.Pool, error) {
	connStr := c.DB.URL
	if connStr == "" {
		return nil, fmt.Errorf("DB_CONN_STR is not set. Check env file")
	}

	db, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	logger.Info("Connected to PostgreSQL")
	return db, nil
}
