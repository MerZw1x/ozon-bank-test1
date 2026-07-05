package postgre

import (
	"backend/src/internal/domain"
	"backend/src/internal/model"
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LinksRepository struct {
	pool *pgxpool.Pool
}

func NewLinksRepository(pool *pgxpool.Pool) *LinksRepository {
	return &LinksRepository{pool: pool}
}

func (r *LinksRepository) Get(ctx context.Context, shortLink string) (domain.Link, error) {
	const q = `SELECT id, original_link, short_link, created_at FROM links WHERE short_link = $1`

	var link model.Link
	err := r.pool.QueryRow(ctx, q, shortLink).Scan(
		&link.ID,
		&link.OriginalLink,
		&link.ShortLink,
		&link.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Link{}, model.ErrNotFound
		}
		return domain.Link{}, err
	}
	return link.ToDomain(), nil
}

func (r *LinksRepository) Save(ctx context.Context, originalLink, shortLink string) (domain.Link, error) {
	const q = `
		INSERT INTO links (original_link, short_link)
		VALUES ($1, $2)
		ON CONFLICT (original_link) DO UPDATE SET original_link = EXCLUDED.original_link
		RETURNING id, original_link, short_link, created_at`

	var link model.Link
	err := r.pool.QueryRow(ctx, q, originalLink, shortLink).Scan(
		&link.ID,
		&link.OriginalLink,
		&link.ShortLink,
		&link.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return domain.Link{}, model.ErrLinkCollision
		}
		return domain.Link{}, err
	}
	return link.ToDomain(), nil
}
