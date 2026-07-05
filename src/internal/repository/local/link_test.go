package local

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinksRepository_Save(t *testing.T) {
	repo := NewLinksRepository()
	ctx := context.Background()

	original := "https://example.com"
	short := "abc123"

	link, err := repo.Save(ctx, original, short)
	require.NoError(t, err)
	require.NotNil(t, link)
	assert.Equal(t, original, link.OriginalLink)
	assert.Equal(t, short, link.ShortLink)
	assert.NotEmpty(t, link.Id)
	assert.False(t, link.CreatedAt.IsZero())

	saved, err := repo.Get(ctx, short)
	require.NoError(t, err)
	require.NotNil(t, saved)
	assert.Equal(t, link.Id, saved.Id)
	assert.Equal(t, original, saved.OriginalLink)
	assert.Equal(t, short, saved.ShortLink)
	assert.Equal(t, link.CreatedAt, saved.CreatedAt)
}

func TestLinksRepository_Save_DuplicateURL(t *testing.T) {
	repo := NewLinksRepository()
	ctx := context.Background()
	original := "https://example.com"
	short1 := "abc123"
	short2 := "xyz789"

	first, err := repo.Save(ctx, original, short1)
	require.NoError(t, err)
	require.NotNil(t, first)

	second, err := repo.Save(ctx, original, short2)
	require.NoError(t, err)
	require.NotNil(t, second)

	assert.Equal(t, first.Id, second.Id)
	assert.Equal(t, first.ShortLink, second.ShortLink)
	assert.Equal(t, first.CreatedAt, second.CreatedAt)

	_, ok := repo.origMap[original]
	assert.True(t, ok)
	assert.Len(t, repo.shortMap, 1)
	assert.Len(t, repo.origMap, 1)
}

func TestLinksRepository_Save_Collision(t *testing.T) {
	repo := NewLinksRepository()
	ctx := context.Background()

	short := "abc123"
	_, err := repo.Save(ctx, "https://google.com", short)
	require.NoError(t, err)

	_, err = repo.Save(ctx, "https://yandex.ru", short)
	assert.ErrorIs(t, err, ErrLinkCollision)

	assert.Len(t, repo.shortMap, 1)
	_, exists := repo.shortMap[short]
	assert.True(t, exists)
}

func TestLinksRepository_Get_NotFound(t *testing.T) {
	repo := NewLinksRepository()
	ctx := context.Background()

	found, err := repo.Get(ctx, "nonexistent")
	require.NoError(t, err)
	assert.Nil(t, found)
}

func TestLinksRepository_Get_EmptyShortLink(t *testing.T) {
	repo := NewLinksRepository()
	ctx := context.Background()

	found, err := repo.Get(ctx, "")
	assert.Error(t, err)
	assert.Nil(t, found)
	assert.EqualError(t, err, ErrShortLinkIsEmpty.Error())
}

func TestLinksRepository_Save_EmptyInput(t *testing.T) {
	repo := NewLinksRepository()
	ctx := context.Background()

	original := "https://example.com"
	short := "abc123"

	link, err := repo.Save(ctx, "", short)
	assert.Error(t, err)
	assert.Nil(t, link)
	assert.EqualError(t, err, ErrOrigLinkIsEmpty.Error())

	link, err = repo.Save(ctx, original, "")
	assert.Error(t, err)
	assert.Nil(t, link)
	assert.EqualError(t, err, ErrShortLinkIsEmpty.Error())
}

func TestLinksRepository_Concurrency(t *testing.T) {
	repo := NewLinksRepository()
	ctx := context.Background()

	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)

	// Каждая горутина сохраняет свой уникальный URL
	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()

			url := "https://example.com/" + string(rune('a'+idx))
			short := "code" + string(rune('a'+idx))

			link, err := repo.Save(ctx, url, short)
			assert.NoError(t, err)
			assert.NotNil(t, link)
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

			short := "code" + string(rune('a'+idx))
			link, err := repo.Get(ctx, short)

			assert.NoError(t, err)
			assert.NotNil(t, link)
		}(i)
	}
	wg2.Wait()
}
