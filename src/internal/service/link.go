package service

import (
	"backend/src/internal/domain"
	"backend/src/internal/repository"
	"backend/src/internal/repository/local"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"strings"
)

const (
	alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_"
	base     = uint64(len(alphabet))
)

type LinksService struct {
	repo repository.ILinksRepository
}

func NewLinksService(repo repository.ILinksRepository) *LinksService {
	return &LinksService{
		repo: repo,
	}
}

func (s *LinksService) SaveLink(originalLink string) (*domain.Link, error) {
	if originalLink == "" {
		return nil, errors.New("original link can not be empty")
	}

	var link *domain.Link
	var err error

	for i := 0; i < 3; i++ {
		shortLink := generateShortLink(originalLink)
		link, err = s.repo.Save(context.Background(), originalLink, shortLink)
		if err != nil {
			if errors.Is(err, local.ErrLinkCollision) {
				continue
			}
			return nil, err
		}
		break
	}

	return link, nil
}

func (s *LinksService) GetLink(shortLink string) (*domain.Link, error) {
	if shortLink == "" {
		return nil, errors.New("short link can not be empty")
	}

	link, err := s.repo.Get(context.Background(), shortLink)
	if err != nil {
		return nil, err
	}

	return link, nil
}

func generateShortLink(originalLink string) string {
	if originalLink == "" {
		return ""
	}

	hash := sha256.Sum256([]byte(originalLink))

	num := binary.BigEndian.Uint64(hash[:8])

	code := encodeBase62(num)

	if len(code) < 10 {
		code = strings.Repeat("0", 10-len(code)) + code
	} else if len(code) > 10 {
		maxVal := uint64(1)
		for i := 0; i < 10; i++ {
			maxVal *= base
		}

		num = num % maxVal
		code = encodeBase62(num)
		if len(code) < 10 {
			code += strings.Repeat("0", 10-len(code))
		}
	}

	return code
}

func encodeBase62(num uint64) string {
	if num == 0 {
		return string(alphabet[0])
	}

	var result []byte
	for num > 0 {
		result = append(result, alphabet[num%base])
		num /= base
	}

	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}
