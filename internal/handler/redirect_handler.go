package handler

import (
	"context"
	"net/http"
	"shorter/internal/events"
	"shorter/internal/metrics"
	"shorter/internal/producer"
	"shorter/internal/repository"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type RedirectHandler struct {
	repo     *repository.LinkRepository
	producer *producer.KafkaProducer
	logger   *zap.Logger
}

func NewRedirectHandler(
	repo *repository.LinkRepository,
	producer *producer.KafkaProducer,
	logger *zap.Logger,
) *RedirectHandler {
	return &RedirectHandler{
		repo:     repo,
		producer: producer,
		logger:   logger,
	}
}

func (rh *RedirectHandler) Handle(w http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "alias")
	if alias == "" {
		metrics.RedirectsErrorTotal.Inc()
		http.Error(w, "alias is required", http.StatusBadRequest)
		return
	}

	link, err := rh.repo.GetByAlias(r.Context(), alias)
	if err != nil {
		metrics.RedirectsErrorTotal.Inc()
		rh.logger.Error("db error", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if link == nil {
		metrics.RedirectsErrorTotal.Inc()
		http.Error(w, "internal error", http.StatusNotFound)
		return
	}

	// update ckicks count
	if err := rh.repo.IncClickCount(r.Context(), alias); err != nil {
		rh.logger.Error("failed to increment click count", zap.Error(err))
	}

	metrics.RedirectsTotal.WithLabelValues(alias).Inc()

	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	userAgent := r.Header.Get("User-Agent")
	referer := r.Header.Get("Referer")

	// create event
	event := &events.ClickEvent{
		Alias:     alias,
		Timestamp: time.Now().UTC(),
		IP:        ip,
		UserAgent: userAgent,
		Referer:   referer,
	}

	// send event to kafka
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := rh.producer.SendClick(ctx, event); err != nil {
			rh.logger.Error("failed to send click event", zap.Error(err))
		}
	}()

	http.Redirect(w,r , link.OriginalUrl, http.StatusFound)
}
