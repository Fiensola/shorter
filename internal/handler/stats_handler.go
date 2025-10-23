package handler

import (
	"encoding/json"
	"net/http"
	"shorter/internal/repository"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type StatsHandler struct {
	repo   repository.AnalyticsRepository
	logger *zap.Logger
}

func NewStatsHandler(repo repository.AnalyticsRepository, logger *zap.Logger) *StatsHandler {
	return &StatsHandler{
		repo: repo,
		logger: logger,
	}
}

func (s *StatsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "alias")
	if alias == "" {
		http.Error(w, "alias is required", http.StatusBadRequest)
		return
	}

	stats, err := s.repo.GetStats(r.Context(), alias)
	if err != nil {
		s.logger.Error("fail to get stats", zap.Error(err), zap.String("alias", alias))
		http.Error(w, "internal error", http.StatusInternalServerError)
        return
	}

	if stats == nil {
        http.Error(w, "stats not found", http.StatusNotFound)
        return
    }

	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stats)
}
