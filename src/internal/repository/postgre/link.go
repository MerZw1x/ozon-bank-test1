package postgre

import (
	"backend/src/internal/domain"
	"backend/src/internal/repository"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LinksRepository struct {
	conn *pgxpool.Pool
}

func NewLinksRepository(conn *pgxpool.Pool) repository.ILinksRepository {
	return &LinksRepository{
		conn: conn,
	}
}

func (r *LinksRepository) Get(ctx context.Context, shortLink string) (*domain.Link, error) {

}

func (r *LinksRepository) Save(ctx context.Context, originalLink, shortLink string) (*domain.Link, error) {

}
