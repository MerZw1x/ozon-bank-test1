package repository

import (
	"backend/src/internal/domain"
	"context"
)

type ILinksRepository interface {
	Get(ctx context.Context, shortLink string) (*domain.Link, error)
	Save(ctx context.Context, originalLink, shortLink string) (*domain.Link, error)
}
