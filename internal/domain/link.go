package domain

import (
	"time"

	"github.com/google/uuid"
)

type Link struct {
	ID           uuid.UUID
	OriginalLink string
	ShortLink    string
	CreatedAt    time.Time
}
