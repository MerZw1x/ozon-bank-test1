package repository

import (
	"backend/src/internal/domain"
	"context"
)

type ILinksRepository interface {
	Get(ctx context.Context, short_link string) (*domain.Link, error)
	Save(ctx context.Context, original_link, short_link string) (*domain.Link, error)
}
