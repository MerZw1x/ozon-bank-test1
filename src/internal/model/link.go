package model

import (
	"time"

	"github.com/google/uuid"
)

type Link struct {
	Id           uuid.UUID `json:"id"`
	OriginalLink string    `json:"original_link"`
	ShortLink    string    `json:"short_link"`
	CreatedAt    time.Time `json:"created_at"`
}
