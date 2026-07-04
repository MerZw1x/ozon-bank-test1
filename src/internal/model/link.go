package model

import (
	"backend/src/internal/domain"
	"time"

	"github.com/google/uuid"
)

type Link struct {
	Id           uuid.UUID `json:"id"`
	OriginalLink string    `json:"original_link"`
	ShortLink    string    `json:"short_link"`
	CreatedAt    time.Time `json:"created_at"`
}

func (l *Link) ToDomain() *domain.Link {
	return &domain.Link{
		Id:           l.Id,
		OriginalLink: l.OriginalLink,
		ShortLink:    l.ShortLink,
		CreatedAt:    l.CreatedAt,
	}
}
