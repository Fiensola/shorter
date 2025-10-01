package enricher

import (
	"context"
)

type Enricher interface {
	Enrich(ctx context.Context, event *ClickTask) (*EnrichedClick, error)
}

type ClickTask struct {
	Alias     string `json:"alias"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	Referer   string `json:"referer"`
	Timestamp string `json:"timestamp"`
}

type EnrichedClick struct {
	Alias     string  `db:"alias"`
	IP        string  `db:"ip"`
	Country   *string `db:"country"`
	City      *string `db:"city"`
	Device    string  `db:"device_type"`
	OS        string  `db:"os"`
	Browser   string  `db:"browser"`
	Referer   *string `db:"referer"`
	Timestamp string  `db:"timestamp"`
}
