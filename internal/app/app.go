package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"shorter/internal/config"
	"shorter/internal/consumer"
	"shorter/internal/enricher"
	"shorter/internal/handler"
	"shorter/internal/logger"
	"shorter/internal/metrics"
	"shorter/internal/producer"
	"shorter/internal/repository"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type App struct {
	Router        *chi.Mux
	Server        *http.Server
	Logger        *zap.Logger
	Db            *pgxpool.Pool
	KafkaProducer *producer.KafkaProducer
	KafkaConsumer *consumer.KafkaConsumer
}

func NewApp() *App {
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

	// metrics
	metrics.Register()

	// repos
	linkRepo := repository.NewLinkRepository(db)
	analyticsRepo := repository.NewAnalyticsRepository(db)

	// kafka
	kafkaProducer := producer.NewKafkaProducer(cfg.Kafka.Brokers, "click_events", logger)
	enricher := enricher.NewIpGeoEnricher(
		enricher.NewGeoClient(cfg.External.GeoAPIkey),
		enricher.NewDeviceParser(),
	)
	kafkaConsumer := consumer.NewKafkaConsumer(
		cfg.Kafka.Brokers,
		"click_events",
		"shorter-consumer-group",
		enricher,
		analyticsRepo,
		logger,
	)

	// router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})

	statsHandler := handler.NewStatsHandler(analyticsRepo, logger)
	r.Get("/api/v1/stats/{alias}", statsHandler.Handle)

	shorterHandler := handler.NewShorterHandler(linkRepo, logger, cfg)
	r.Post("/api/v1/shorter", shorterHandler.Handle)

	redirectHandler := handler.NewRedirectHandler(linkRepo, kafkaProducer, logger)
	r.Get("/{alias}", redirectHandler.Handle)

	// http
	srv := &http.Server{
		Addr:         cfg.Server.Host + ":" + fmt.Sprint(cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &App{
		Router:        r,
		Server:        srv,
		Logger:        logger,
		Db:            db,
		KafkaProducer: kafkaProducer,
		KafkaConsumer: kafkaConsumer,
	}
}

func (a *App) StartHttpServer() error {
	a.Logger.Info("starting HTTP server at", zap.String("addr", a.Server.Addr))
	return a.Server.ListenAndServe()
}

func (a *App) ShutdownHttpServer(ctx context.Context) error {
	return a.Server.Shutdown(ctx)
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
