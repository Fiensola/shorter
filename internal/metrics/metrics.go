package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	RedirectsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "shorter",
			Subsystem: "http",
			Name: "redirects_total",
			Help: "Total number of redirects by alias",
		},
		[]string{"alias"},
	)

	RedirectsErrorTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "shorter",
			Subsystem: "http",
			Name: "redirect_errors_total",
			Help: "Total number of redirect errors",
		},
	)

	EnrichDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "shorter",
			Subsystem: "enricher",
			Name: "duration_seconds",
			Help: "Enrichment duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)

	EventsProcessed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "shorter",
			Subsystem: "consumer",
			Name: "events_processed_total",
			Help: "Total number of events prcessed by consumer",
		},
	)
)

func Register() {
	prometheus.MustRegister(RedirectsTotal)
	prometheus.MustRegister(RedirectsErrorTotal)
	prometheus.MustRegister(EnrichDuration)
	prometheus.MustRegister(EventsProcessed)
}