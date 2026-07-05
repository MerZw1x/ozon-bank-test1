package model

import (
	"backend/internal/domain"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotFound      = errors.New("link not found")
	ErrLinkCollision = errors.New("short link collision")
)

type Link struct {
	ID           uuid.UUID `db:"id"`
	OriginalLink string    `db:"original_link"`
	ShortLink    string    `db:"short_link"`
	CreatedAt    time.Time `db:"created_at"`
}

func (l Link) ToDomain() domain.Link {
	return domain.Link{
		ID:           l.ID,
		OriginalLink: l.OriginalLink,
		ShortLink:    l.ShortLink,
		CreatedAt:    l.CreatedAt,
	}
}
