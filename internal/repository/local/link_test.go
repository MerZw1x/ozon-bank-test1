package local

import (
	"backend/internal/model"
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinksRepository_Save(t *testing.T) {
	repo := NewLinksRepository()
	ctx := context.Background()

	original := "https://example.com"
	short := "abc1234567"

	link, err := repo.Save(ctx, original, short)
	require.NoError(t, err)
	assert.Equal(t, original, link.OriginalLink)
	assert.Equal(t, short, link.ShortLink)
	assert.NotEmpty(t, link.ID)
	assert.False(t, link.CreatedAt.IsZero())

	saved, err := repo.Get(ctx, short)
	require.NoError(t, err)
	assert.Equal(t, link, saved)
}

func TestLinksRepository_Save_DuplicateURL(t *testing.T) {
	repo := NewLinksRepository()
	ctx := context.Background()
	original := "https://example.com"

	first, err := repo.Save(ctx, original, "abc1234567")
	require.NoError(t, err)

	second, err := repo.Save(ctx, original, "xyz9876543")
	require.NoError(t, err)

	assert.Equal(t, first, second)
	assert.Len(t, repo.shortMap, 1)
	assert.Len(t, repo.origMap, 1)
}

func TestLinksRepository_Save_Collision(t *testing.T) {
	repo := NewLinksRepository()
	ctx := context.Background()

	short := "abc1234567"
	_, err := repo.Save(ctx, "https://google.com", short)
	require.NoError(t, err)

	_, err = repo.Save(ctx, "https://yandex.ru", short)
	assert.ErrorIs(t, err, model.ErrLinkCollision)

	assert.Len(t, repo.shortMap, 1)
}

func TestLinksRepository_Get_NotFound(t *testing.T) {
	repo := NewLinksRepository()

	_, err := repo.Get(context.Background(), "nonexistent")
	assert.ErrorIs(t, err, model.ErrNotFound)
}

func TestLinksRepository_Concurrency(t *testing.T) {
	repo := NewLinksRepository()
	ctx := context.Background()

	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()

			url := fmt.Sprintf("https://example.com/%d", idx)
			short := fmt.Sprintf("code%06d", idx)

			_, err := repo.Save(ctx, url, short)
			assert.NoError(t, err)
		}(i)
	}
	wg.Wait()

	assert.Len(t, repo.shortMap, goroutines)
	assert.Len(t, repo.origMap, goroutines)

	var wg2 sync.WaitGroup
	wg2.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg2.Done()

			short := fmt.Sprintf("code%06d", idx)
			_, err := repo.Get(ctx, short)
			assert.NoError(t, err)
		}(i)
	}
	wg2.Wait()
}
