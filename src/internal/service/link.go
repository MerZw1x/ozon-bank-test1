package service

import (
	"backend/src/internal/domain"
	"backend/src/internal/repository"
	"context"
	"crypto/sha256"
	"errors"
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

	shortLink := generateShortLink(originalLink)

	link, err := s.repo.Save(context.Background(), originalLink, shortLink)
	if err != nil {
		return nil, err
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
	return string(hash[:10])
}
