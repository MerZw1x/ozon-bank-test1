package postgre

import (
	"backend/src/internal/domain"
	"backend/src/internal/model"
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LinksRepository struct {
	conn *pgxpool.Pool
}

func NewLinksRepository(conn *pgxpool.Pool) *LinksRepository {
	return &LinksRepository{
		conn: conn,
	}
}

func (r *LinksRepository) Get(ctx context.Context, shortLink string) (*domain.Link, error) {
	sqlStr := `SELECT id, original_link, short_link, created_at FROM links WHERE short_link = $1`

	link := &model.Link{}
	err := r.conn.QueryRow(ctx, sqlStr, shortLink).Scan(
		&link.Id,
		&link.OriginalLink,
		&link.ShortLink,
		&link.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return link.ToDomain(), nil
}

func (r *LinksRepository) Save(ctx context.Context, originalLink, shortLink string) (*domain.Link, error) {
	sqlStr := `INSERT INTO link (original_link, short_link)
				VALUES ($1, $2)
				RETURNING id, original_link, short_link, created_at`

	link := &model.Link{}
	err := r.conn.QueryRow(ctx, sqlStr, originalLink, shortLink).Scan(
		&link.Id,
		&link.OriginalLink,
		&link.ShortLink,
		&link.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return link.ToDomain(), nil
}
