package service

import (
	"backend/src/internal/domain"
	"backend/src/internal/model"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/fnv"
)

const (
	alphabet   = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_"
	base       = uint64(len(alphabet))
	shortLen   = 10
	maxRetries = 5
)

type LinksRepository interface {
	Get(ctx context.Context, shortLink string) (domain.Link, error)
	Save(ctx context.Context, originalLink, shortLink string) (domain.Link, error)
}

type LinksService struct {
	repo LinksRepository
}

func NewLinksService(repo LinksRepository) *LinksService {
	return &LinksService{repo: repo}
}

func (s *LinksService) Shorten(ctx context.Context, originalLink string) (domain.Link, error) {
	for i := 0; i < maxRetries; i++ {
		short := generateShortLink(originalLink, i)
		link, err := s.repo.Save(ctx, originalLink, short)
		if err == nil {
			return link, nil
		}
		if !errors.Is(err, model.ErrLinkCollision) {
			return domain.Link{}, err
		}
	}
	return domain.Link{}, fmt.Errorf("failed to generate unique short link after %d attempts", maxRetries)
}

func (s *LinksService) GetOriginal(ctx context.Context, shortLink string) (domain.Link, error) {
	return s.repo.Get(ctx, shortLink)
}

func generateShortLink(originalLink string, salt int) string {
	h := fnv.New64a()
	h.Write([]byte(originalLink))
	if salt > 0 {
		var buf [8]byte
		binary.BigEndian.PutUint64(buf[:], uint64(salt))
		h.Write(buf[:])
	}
	return encode(h.Sum64())
}

func encode(n uint64) string {
	var b [shortLen]byte
	for i := shortLen - 1; i >= 0; i-- {
		b[i] = alphabet[n%base]
		n /= base
	}
	return string(b[:])
}
