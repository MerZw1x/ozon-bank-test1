package model

import (
	"backend/src/internal/domain"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotFound      = errors.New("link not found")
	ErrLinkCollision = errors.New("short link collision")
)

type Link struct {
	ID           uuid.UUID
	OriginalLink string
	ShortLink    string
	CreatedAt    time.Time
}

func (l Link) ToDomain() domain.Link {
	return domain.Link{
		ID:           l.ID,
		OriginalLink: l.OriginalLink,
		ShortLink:    l.ShortLink,
		CreatedAt:    l.CreatedAt,
	}
}
