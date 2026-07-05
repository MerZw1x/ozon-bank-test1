package service

import (
	"backend/src/internal/domain"
	"backend/src/internal/repository/local"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRepository struct {
	saveFunc func(ctx context.Context, originalLink, shortLink string) (*domain.Link, error)
	getFunc  func(ctx context.Context, shortLink string) (*domain.Link, error)
}

func (m *mockRepository) Save(ctx context.Context, originalLink, shortLink string) (*domain.Link, error) {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, originalLink, shortLink)
	}
	return nil, nil
}

func (m *mockRepository) Get(ctx context.Context, shortLink string) (*domain.Link, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, shortLink)
	}
	return nil, nil
}

func (m *mockRepository) Close() error {
	return nil
}

func TestLinksService_SaveLink_Success(t *testing.T) {
	original := "https://example.com"
	expectedShort := "abc123def0"
	expectedLink := &domain.Link{
		Id:           uuid.New(),
		OriginalLink: original,
		ShortLink:    expectedShort,
		CreatedAt:    time.Now(),
	}

	mock := &mockRepository{
		saveFunc: func(ctx context.Context, orig, short string) (*domain.Link, error) {
			return expectedLink, nil
		},
	}

	svc := NewLinksService(mock)
	ctx := context.Background()

	link, err := svc.SaveLink(ctx, original)
	require.NoError(t, err)
	require.NotNil(t, link)
	assert.Equal(t, expectedLink.ShortLink, link.ShortLink)
	assert.Equal(t, expectedLink.OriginalLink, link.OriginalLink)
}

func TestLinksService_SaveLink_Duplicate(t *testing.T) {
	original := "https://example.com"
	existingShort := "abc123"
	existingLink := &domain.Link{
		Id:           uuid.New(),
		OriginalLink: original,
		ShortLink:    existingShort,
		CreatedAt:    time.Now(),
	}

	callCount := 0
	mock := &mockRepository{
		saveFunc: func(ctx context.Context, orig, short string) (*domain.Link, error) {
			if callCount == 0 {
				callCount++
				return &domain.Link{
					Id:           uuid.New(),
					OriginalLink: orig,
					ShortLink:    short,
					CreatedAt:    time.Now(),
				}, nil
			}
			return existingLink, nil
		},
	}

	svc := NewLinksService(mock)
	ctx := context.Background()

	first, err := svc.SaveLink(ctx, original)
	require.NoError(t, err)
	require.NotNil(t, first)

	second, err := svc.SaveLink(ctx, original)
	require.NoError(t, err)
	require.NotNil(t, second)

	assert.Equal(t, first.ShortLink, second.ShortLink)
}

func TestLinksService_SaveLink_CollisionRetry(t *testing.T) {
	original := "https://example.com"
	mock := &mockRepository{
		saveFunc: func(ctx context.Context, orig, short string) (*domain.Link, error) {
			if short == "abc123def0" {
				return nil, local.ErrLinkCollision
			}

			return &domain.Link{
				Id:           uuid.New(),
				OriginalLink: orig,
				ShortLink:    short,
				CreatedAt:    time.Now(),
			}, nil
		},
	}

	svc := NewLinksService(mock)
	ctx := context.Background()

	link, err := svc.SaveLink(ctx, original)
	require.NoError(t, err)
	require.NotNil(t, link)
	assert.NotEmpty(t, link.ShortLink)
}

func TestLinksService_SaveLink_EmptyOriginal(t *testing.T) {
	svc := NewLinksService(&mockRepository{})
	ctx := context.Background()

	link, err := svc.SaveLink(ctx, "")
	assert.Error(t, err)
	assert.Nil(t, link)
	assert.EqualError(t, err, "original link can not be empty")
}
func TestLinksService_GetLink_EmptyShort(t *testing.T) {
	svc := NewLinksService(&mockRepository{})
	ctx := context.Background()

	link, err := svc.GetLink(ctx, "")
	assert.Error(t, err)
	assert.Nil(t, link)
	assert.EqualError(t, err, "short link can not be empty")
}

func TestLinksService_GetLink_Success(t *testing.T) {
	short := "abc123"
	expectedLink := &domain.Link{
		Id:           uuid.New(),
		OriginalLink: "https://example.com",
		ShortLink:    short,
		CreatedAt:    time.Now(),
	}

	mock := &mockRepository{
		getFunc: func(ctx context.Context, s string) (*domain.Link, error) {
			assert.Equal(t, short, s)
			return expectedLink, nil
		},
	}

	svc := NewLinksService(mock)
	ctx := context.Background()

	link, err := svc.GetLink(ctx, short)
	require.NoError(t, err)
	require.NotNil(t, link)
	assert.Equal(t, expectedLink.ShortLink, link.ShortLink)
	assert.Equal(t, expectedLink.OriginalLink, link.OriginalLink)
}

func TestLinksService_GetLink_NotFound(t *testing.T) {
	mock := &mockRepository{
		getFunc: func(ctx context.Context, s string) (*domain.Link, error) {
			return nil, nil
		},
	}

	svc := NewLinksService(mock)
	ctx := context.Background()

	link, err := svc.GetLink(ctx, "nonexistent")
	require.NoError(t, err)
	assert.Nil(t, link)
}
