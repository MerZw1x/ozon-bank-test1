package local

import (
	"backend/src/internal/domain"
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type LinksRepository struct {
	mu       sync.RWMutex
	shortMap map[string]*domain.Link
	origMap  map[string]string
}

func NewLinksRepository() *LinksRepository {
	return &LinksRepository{
		shortMap: make(map[string]*domain.Link),
		origMap:  make(map[string]string),
	}
}

func (r *LinksRepository) Get(ctx context.Context, shortLink string) (*domain.Link, error) {
	if shortLink == "" {
		return nil, errors.New("short link can not be empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	if value, ok := r.shortMap[shortLink]; !ok {
		return nil, nil
	} else {
		return value, nil
	}
}

func (r *LinksRepository) Save(ctx context.Context, originalLink, shortLink string) (*domain.Link, error) {
	if shortLink == "" {
		return nil, errors.New("short link can not be empty")
	}

	if originalLink == "" {
		return nil, errors.New("original link can not be empty")
	}

	if value, ok := r.shortMap[shortLink]; ok {
		return value, nil
	}

	link := &domain.Link{
		Id:           uuid.New(),
		OriginalLink: originalLink,
		ShortLink:    shortLink,
		CreatedAt:    time.Now(),
	}

	r.shortMap[shortLink] = link
	r.origMap[originalLink] = shortLink
	return link, nil
}
