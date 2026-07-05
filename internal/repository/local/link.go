package local

import (
	"backend/internal/domain"
	"backend/internal/model"
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type LinksRepository struct {
	mu       sync.RWMutex
	shortMap map[string]model.Link
	origMap  map[string]string
}

func NewLinksRepository() *LinksRepository {
	return &LinksRepository{
		shortMap: make(map[string]model.Link),
		origMap:  make(map[string]string),
	}
}

func (r *LinksRepository) Ping(_ context.Context) error {
	return nil
}

func (r *LinksRepository) Get(_ context.Context, shortLink string) (domain.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	link, ok := r.shortMap[shortLink]
	if !ok {
		return domain.Link{}, model.ErrNotFound
	}
	return link.ToDomain(), nil
}

func (r *LinksRepository) Save(_ context.Context, originalLink, shortLink string) (domain.Link, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if existingShort, ok := r.origMap[originalLink]; ok {
		return r.shortMap[existingShort].ToDomain(), nil
	}

	if _, ok := r.shortMap[shortLink]; ok {
		return domain.Link{}, model.ErrLinkCollision
	}

	link := model.Link{
		ID:           uuid.New(),
		OriginalLink: originalLink,
		ShortLink:    shortLink,
		CreatedAt:    time.Now(),
	}

	r.shortMap[shortLink] = link
	r.origMap[originalLink] = shortLink
	return link.ToDomain(), nil
}
