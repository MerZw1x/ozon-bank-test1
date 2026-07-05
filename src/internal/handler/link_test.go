package handler

import (
	"backend/src/internal/domain"
	"backend/src/internal/dto"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockService struct {
	saveFunc func(ctx context.Context, originalURL string) (*domain.Link, error)
	getFunc  func(ctx context.Context, shortLink string) (*domain.Link, error)
}

func (m *mockService) SaveLink(ctx context.Context, originalURL string) (*domain.Link, error) {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, originalURL)
	}
	return nil, nil
}

func (m *mockService) GetLink(ctx context.Context, shortLink string) (*domain.Link, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, shortLink)
	}
	return nil, nil
}

func setupApp(mock *mockService) *fiber.App {
	app := fiber.New()
	handler := &LinksHandler{service: mock}
	app.Post("/shorten", handler.Shorten)
	app.Get("/:shortLink", handler.Redirect)
	return app
}

func TestLinksHandler_Shorten_Success(t *testing.T) {
	originalURL := "https://example.com"
	shortLink := "abc123"
	expectedLink := &domain.Link{
		Id:           uuid.New(),
		OriginalLink: originalURL,
		ShortLink:    shortLink,
		CreatedAt:    time.Now(),
	}

	mock := &mockService{
		saveFunc: func(ctx context.Context, url string) (*domain.Link, error) {
			assert.Equal(t, originalURL, url)
			return expectedLink, nil
		},
	}

	app := setupApp(mock)

	body := dto.SaveLinkRequest{
		URL: originalURL,
	}

	req := httptest.NewRequest("POST", "/shorten", bytes.NewReader([]byte(body.URL)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, shortLink, response["short_link"])
	assert.Empty(t, response["error"])
}

func TestLinksHandler_Shorten_InvalidJSON(t *testing.T) {
	mock := &mockService{}
	app := setupApp(mock)

	req := httptest.NewRequest("POST", "/shorten", bytes.NewReader([]byte(`{"url": "https://example.com"`)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "invalid request body", response["error"])
}

func TestLinksHandler_Shorten_EmptyURL(t *testing.T) {
	mock := &mockService{}
	app := setupApp(mock)

	reqBody := map[string]string{"url": ""}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/shorten", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "url is required", response["error"])
}

func TestLinksHandler_Shorten_InvalidURL(t *testing.T) {
	mock := &mockService{}
	app := setupApp(mock)

	reqBody := map[string]string{"url": "example.com"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/shorten", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "invalid URL format", response["error"])
}

func TestLinksHandler_Shorten_ServiceError(t *testing.T) {
	mock := &mockService{
		saveFunc: func(ctx context.Context, url string) (*domain.Link, error) {
			return nil, errors.New("internal error")
		},
	}
	app := setupApp(mock)

	reqBody := map[string]string{"url": "https://example.com"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/shorten", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "failed to shorten URL", response["error"])
}

func TestLinksHandler_Redirect_Success(t *testing.T) {
	shortLink := "abc123"
	originalURL := "https://example.com"
	expectedLink := &domain.Link{
		Id:           uuid.New(),
		OriginalLink: originalURL,
		ShortLink:    shortLink,
		CreatedAt:    time.Now(),
	}

	mock := &mockService{
		getFunc: func(ctx context.Context, s string) (*domain.Link, error) {
			assert.Equal(t, shortLink, s)
			return expectedLink, nil
		},
	}
	app := setupApp(mock)

	req := httptest.NewRequest("GET", "/"+shortLink, nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, originalURL, response["original_url"])
	assert.Empty(t, response["error"])
}

func TestLinksHandler_Redirect_NotFound(t *testing.T) {
	mock := &mockService{
		getFunc: func(ctx context.Context, s string) (*domain.Link, error) {
			return nil, nil
		},
	}
	app := setupApp(mock)

	req := httptest.NewRequest("GET", "/nonexistent", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "short link not found", response["error"])
}

func TestLinksHandler_Redirect_EmptyShort(t *testing.T) {
	mock := &mockService{}
	app := setupApp(mock)

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "short link is required", response["error"])
}

func TestLinksHandler_Redirect_ServiceError(t *testing.T) {
	mock := &mockService{
		getFunc: func(ctx context.Context, s string) (*domain.Link, error) {
			return nil, errors.New("db error")
		},
	}
	app := setupApp(mock)

	req := httptest.NewRequest("GET", "/abc123", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "failed to retrieve original URL", response["error"])
}
