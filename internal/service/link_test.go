package service

import (
	"backend/internal/domain"
	"backend/internal/model"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRepository struct {
	saveFunc func(ctx context.Context, originalLink, shortLink string) (domain.Link, error)
	getFunc  func(ctx context.Context, shortLink string) (domain.Link, error)
}

func (m *mockRepository) Save(ctx context.Context, originalLink, shortLink string) (domain.Link, error) {
	return m.saveFunc(ctx, originalLink, shortLink)
}

func (m *mockRepository) Get(ctx context.Context, shortLink string) (domain.Link, error) {
	return m.getFunc(ctx, shortLink)
}

func (m *mockRepository) Ping(_ context.Context) error {
	return nil
}

func TestShorten_Success(t *testing.T) {
	original := "https://example.com"

	mock := &mockRepository{
		saveFunc: func(_ context.Context, orig, short string) (domain.Link, error) {
			assert.Len(t, short, shortLen)
			return domain.Link{ID: uuid.New(), OriginalLink: orig, ShortLink: short}, nil
		},
	}

	link, err := NewLinksService(mock).Shorten(context.Background(), original)
	require.NoError(t, err)
	assert.Equal(t, original, link.OriginalLink)
	assert.Len(t, link.ShortLink, shortLen)
}

func TestShorten_Deterministic(t *testing.T) {
	assert.Equal(t, generateShortLink("https://example.com", 0), generateShortLink("https://example.com", 0))
	assert.NotEqual(t, generateShortLink("https://example.com", 0), generateShortLink("https://example.com", 1))
	assert.NotEqual(t, generateShortLink("https://a.com", 0), generateShortLink("https://b.com", 0))
}

func TestShorten_CollisionRetry(t *testing.T) {
	original := "https://example.com"
	firstShort := generateShortLink(original, 0)

	calls := 0
	mock := &mockRepository{
		saveFunc: func(_ context.Context, orig, short string) (domain.Link, error) {
			calls++
			if short == firstShort {
				return domain.Link{}, model.ErrLinkCollision
			}
			return domain.Link{ID: uuid.New(), OriginalLink: orig, ShortLink: short}, nil
		},
	}

	link, err := NewLinksService(mock).Shorten(context.Background(), original)
	require.NoError(t, err)
	assert.NotEqual(t, firstShort, link.ShortLink)
	assert.Equal(t, 2, calls)
}

func TestShorten_CollisionExhausted(t *testing.T) {
	mock := &mockRepository{
		saveFunc: func(_ context.Context, _, _ string) (domain.Link, error) {
			return domain.Link{}, model.ErrLinkCollision
		},
	}

	_, err := NewLinksService(mock).Shorten(context.Background(), "https://example.com")
	assert.Error(t, err)
}

func TestShorten_RepoError(t *testing.T) {
	repoErr := errors.New("boom")
	mock := &mockRepository{
		saveFunc: func(_ context.Context, _, _ string) (domain.Link, error) {
			return domain.Link{}, repoErr
		},
	}

	_, err := NewLinksService(mock).Shorten(context.Background(), "https://example.com")
	assert.ErrorIs(t, err, repoErr)
}

func TestGetOriginal_Success(t *testing.T) {
	short := "abc1234567"
	expected := domain.Link{ID: uuid.New(), OriginalLink: "https://example.com", ShortLink: short}

	mock := &mockRepository{
		getFunc: func(_ context.Context, s string) (domain.Link, error) {
			assert.Equal(t, short, s)
			return expected, nil
		},
	}

	link, err := NewLinksService(mock).GetOriginal(context.Background(), short)
	require.NoError(t, err)
	assert.Equal(t, expected, link)
}

func TestGetOriginal_NotFound(t *testing.T) {
	mock := &mockRepository{
		getFunc: func(_ context.Context, _ string) (domain.Link, error) {
			return domain.Link{}, model.ErrNotFound
		},
	}

	_, err := NewLinksService(mock).GetOriginal(context.Background(), "nonexistent")
	assert.ErrorIs(t, err, model.ErrNotFound)
}

func TestGenerateShortLink_AlphabetOnly(t *testing.T) {
	for i := 0; i < 100; i++ {
		short := generateShortLink("https://example.com/"+string(rune('a'+i%26)), i)
		assert.Len(t, short, shortLen)
		for _, ch := range short {
			assert.Contains(t, alphabet, string(ch))
		}
	}
}
