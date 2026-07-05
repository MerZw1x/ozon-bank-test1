package local

import (
	"backend/src/internal/domain"
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var ErrLinkCollision error = errors.New("short link collision")
var ErrShortLinkIsEmpty error = errors.New("short link can not be empty")
var ErrOrigLinkIsEmpty error = errors.New("original link can not be empty")

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
		return nil, ErrShortLinkIsEmpty
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	if value, ok := r.shortMap[shortLink]; !ok {
		return nil, nil
	} else {
		return value, nil
	}
}

func (r *LinksRepository) Save(ctx context.Context, originalLink, shortLink string) (*domain.Link, error) {
	if shortLink == "" {
		return nil, ErrShortLinkIsEmpty
	}

	if originalLink == "" {
		return nil, ErrOrigLinkIsEmpty
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if value, ok := r.origMap[originalLink]; ok {
		return r.shortMap[value], nil
	}

	if value, ok := r.shortMap[shortLink]; ok {
		if value.OriginalLink != originalLink {
			return nil, ErrLinkCollision
		}
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
