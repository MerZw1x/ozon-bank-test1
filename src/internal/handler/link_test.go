package handler

import (
	"backend/src/internal/domain"
	"backend/src/internal/dto"
	"backend/src/internal/service"
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

// mockService — реализация интерфейса LinkService для тестов
type mockService struct {
	saveFunc func(ctx context.Context, url string) (*domain.Link, error)
	getFunc  func(ctx context.Context, shortLink string) (*domain.Link, error)
}

func (m *mockService) SaveLink(ctx context.Context, url string) (*domain.Link, error) {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, url)
	}
	return nil, nil
}

func (m *mockService) GetLink(ctx context.Context, shortLink string) (*domain.Link, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, shortLink)
	}
	return nil, nil
}

func setupApp(service service.ILinkService) *fiber.App {
	app := fiber.New()
	handler := NewLinkHandler(service)
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

	reqBody := dto.SaveLinkRequest{URL: originalURL}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/shorten", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, shortLink, response["short_link"])
	assert.Empty(t, response["error"])
}

func TestLinksHandler_Shorten_EmptyURL(t *testing.T) {
	mock := &mockService{}
	app := setupApp(mock)

	reqBody := dto.SaveLinkRequest{URL: ""}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/shorten", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "url is required", response["error"])
}

func TestLinksHandler_Shorten_InvalidURL(t *testing.T) {
	mock := &mockService{}
	app := setupApp(mock)

	reqBody := dto.SaveLinkRequest{URL: "not-a-url"}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/shorten", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "invalid URL format", response["error"])
}

func TestLinksHandler_Shorten_InvalidJSON(t *testing.T) {
	mock := &mockService{}
	app := setupApp(mock)

	// Некорректный JSON
	req := httptest.NewRequest("POST", "/shorten", bytes.NewReader([]byte(`{"url": "broken"`)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "invalid request body", response["error"])
}

func TestLinksHandler_Shorten_ServiceError(t *testing.T) {
	mock := &mockService{
		saveFunc: func(ctx context.Context, url string) (*domain.Link, error) {
			return nil, errors.New("internal error")
		},
	}
	app := setupApp(mock)

	reqBody := dto.SaveLinkRequest{URL: "https://example.com"}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/shorten", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "failed to shorten URL", response["error"])
}

func TestLinksHandler_Redirect_Success(t *testing.T) {
	originalURL := "https://example.com"
	shortLink := "abc123"
	expectedLink := &domain.Link{
		Id:           uuid.New(),
		OriginalLink: originalURL,
		ShortLink:    shortLink,
	}

	mock := &mockService{
		getFunc: func(ctx context.Context, short string) (*domain.Link, error) {
			assert.Equal(t, shortLink, short)
			return expectedLink, nil
		},
	}
	app := setupApp(mock)

	req := httptest.NewRequest("GET", "/"+shortLink, nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, originalURL, response["original_url"])
	assert.Empty(t, response["error"])
}

func TestLinksHandler_Redirect_NotFound(t *testing.T) {
	mock := &mockService{
		getFunc: func(ctx context.Context, short string) (*domain.Link, error) {
			return nil, nil
		},
	}
	app := setupApp(mock)

	req := httptest.NewRequest("GET", "/nonexistent", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "short link not found", response["error"])
}

func TestLinksHandler_Redirect_ServiceError(t *testing.T) {
	mock := &mockService{
		getFunc: func(ctx context.Context, short string) (*domain.Link, error) {
			return nil, errors.New("db error")
		},
	}
	app := setupApp(mock)

	req := httptest.NewRequest("GET", "/abc123", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "failed to retrieve original URL", response["error"])
}
