package domain

import (
	"time"

	"github.com/google/uuid"
)

type Link struct {
	Id           uuid.UUID
	OriginalLink string
	ShortLink    string
	CreatedAt    time.Time
}
