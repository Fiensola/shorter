package repository

import (
	"context"
	"shorter/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type LinkRepository interface {
	Create(ctx context.Context, link *model.Link) error
	GetByAlias(ctx context.Context, alias string) (*model.Link, error)
	IncClickCount(ctx context.Context, alias string) error
}

type PgLinkRepository struct {
	db *pgxpool.Pool
}

func NewLinkRepository(db *pgxpool.Pool) *PgLinkRepository {

	return &PgLinkRepository{db: db}
}

func (r *PgLinkRepository) Create(ctx context.Context, link *model.Link) error {
	q := `
		INSERT INTO 
			short_links (alias, original_url, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(ctx, q,
		link.Alias,
		link.OriginalUrl,
		link.ExpiresAt,
	)

	return err
}

func (r *PgLinkRepository) GetByAlias(ctx context.Context, alias string) (*model.Link, error) {
	q := `
		SELECT 
			alias, original_url, expires_at, click_count
		FROM
			short_links
		WHERE alias = $1
	`
	var link model.Link
	err := r.db.QueryRow(ctx, q, alias).Scan(
		&link.Alias,
		&link.OriginalUrl,
		&link.ExpiresAt,
		&link.ClickCount,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}

	return &link, err
}

func (r *PgLinkRepository) IncClickCount(ctx context.Context, alias string) error {
	q := `
		UPDATE short_links
		SET click_count = click_count + 1
		WHERE alias = $1
	`
	_, err := r.db.Exec(ctx, q, alias)
	return err
}
