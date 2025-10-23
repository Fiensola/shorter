package repository

import (
	"context"
	"shorter/internal/enricher"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AnalyticsRepository interface {
	Save(ctx context.Context, click *enricher.EnrichedClick) error
	GetStats(ctx context.Context, alias string) (*Stats, error)
}

type PgAnalyticsRepository struct {
	db *pgxpool.Pool
}

type Stats struct {
	Alias       string         `json:"alias"`
	TotalClicks int            `json:"total_clicks"`
	UniqueIPs   int            `json:"unique_ips"`
	ByCountry   map[string]int `json:"by_country"`
	ByDevice    map[string]int `json:"by_device"`
	ByBrowser   map[string]int `json:"by_browser"`
}

func NewAnalyticsRepository(db *pgxpool.Pool) *PgAnalyticsRepository {
	return &PgAnalyticsRepository{
		db: db,
	}
}

func (r *PgAnalyticsRepository) Save(ctx context.Context, click *enricher.EnrichedClick) error {
	q := `
		INSERT INTO enriched_clicks
		(alias, ip, country, city, device_type, os, browser, referer, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.Exec(ctx, q,
		click.Alias,
		click.IP,
		click.Country,
		click.City,
		click.Device,
		click.OS,
		click.Browser,
		click.Referer,
		click.Timestamp,
	)

	return err
}

func (r *PgAnalyticsRepository) GetStats(ctx context.Context, alias string) (*Stats, error) {
	var stats Stats
	stats.Alias = alias

	q := `
		SELECT 
			country, city, device_type, os
		FROM 
			enriched_clicks
		WHERE 
			alias=$1
	`
	rows, err := r.db.Query(ctx, q, alias)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats.ByCountry = make(map[string]int)
	stats.ByDevice = make(map[string]int)
	stats.ByBrowser = make(map[string]int)

	for rows.Next() {
		// todo count by county, device, browser
	}

	stats.TotalClicks = int(rows.CommandTag().RowsAffected())

	return &stats, nil
}
