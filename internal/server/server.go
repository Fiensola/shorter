package server

import (
	"context"
	"net/http"
	"shorter/internal/config"
	"shorter/internal/handler"
	"shorter/internal/repository"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Server struct {
	router *chi.Mux
	server *http.Server
	logger *zap.Logger
	db     *pgxpool.Pool
}

func NewServer(addr string, logger *zap.Logger, cfg *config.Config, db *pgxpool.Pool) *Server {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	linkRepo := repository.NewLinkRepository(db)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	shorterHandler := handler.NewShorterHandler(linkRepo, logger, cfg)
	r.Post("/api/v1/shorter", shorterHandler.Handle)

	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &Server{
		router: r,
		server: srv,
		logger: logger,
		db: db,
	}
}

func (s *Server) Start() error {
	s.logger.Info("starting HTTP server at", zap.String("addr", s.server.Addr))
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
