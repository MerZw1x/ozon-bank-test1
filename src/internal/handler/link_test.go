package handler

import (
	"backend/src/internal/domain"
	"backend/src/internal/model"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockService struct {
	shortenFunc func(ctx context.Context, url string) (domain.Link, error)
	getFunc     func(ctx context.Context, shortLink string) (domain.Link, error)
}

func (m *mockService) Shorten(ctx context.Context, url string) (domain.Link, error) {
	return m.shortenFunc(ctx, url)
}

func (m *mockService) GetOriginal(ctx context.Context, shortLink string) (domain.Link, error) {
	return m.getFunc(ctx, shortLink)
}

func setupApp(svc LinksService) *fiber.App {
	app := fiber.New()
	NewLinksHandler(svc).Register(app)
	return app
}

func do(t *testing.T, app *fiber.App, method, path string, body []byte) *http.Response {
	t.Helper()
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	return resp
}

func decode(t *testing.T, resp *http.Response) map[string]string {
	t.Helper()
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var out map[string]string
	require.NoError(t, json.Unmarshal(raw, &out))
	return out
}

func TestShortenHandler_Success(t *testing.T) {
	original := "https://example.com"
	short := "abc1234567"

	app := setupApp(&mockService{
		shortenFunc: func(_ context.Context, url string) (domain.Link, error) {
			assert.Equal(t, original, url)
			return domain.Link{ID: uuid.New(), OriginalLink: url, ShortLink: short}, nil
		},
	})

	body, _ := json.Marshal(shortenRequest{URL: original})
	resp := do(t, app, "POST", "/shorten", body)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, short, decode(t, resp)["short_link"])
}

func TestShortenHandler_InvalidJSON(t *testing.T) {
	app := setupApp(&mockService{})
	resp := do(t, app, "POST", "/shorten", []byte(`{"url": "broken"`))

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "invalid request body", decode(t, resp)["error"])
}

func TestShortenHandler_ValidationFailures(t *testing.T) {
	app := setupApp(&mockService{})

	cases := map[string]string{
		"empty":     "",
		"not_a_url": "not-a-url",
	}
	for name, url := range cases {
		t.Run(name, func(t *testing.T) {
			body, _ := json.Marshal(shortenRequest{URL: url})
			resp := do(t, app, "POST", "/shorten", body)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			assert.NotEmpty(t, decode(t, resp)["error"])
		})
	}
}

func TestShortenHandler_ServiceError(t *testing.T) {
	app := setupApp(&mockService{
		shortenFunc: func(_ context.Context, _ string) (domain.Link, error) {
			return domain.Link{}, errors.New("boom")
		},
	})

	body, _ := json.Marshal(shortenRequest{URL: "https://example.com"})
	resp := do(t, app, "POST", "/shorten", body)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, "failed to shorten URL", decode(t, resp)["error"])
}

func TestGetHandler_Success(t *testing.T) {
	original := "https://example.com"
	short := "abc1234567"

	app := setupApp(&mockService{
		getFunc: func(_ context.Context, s string) (domain.Link, error) {
			assert.Equal(t, short, s)
			return domain.Link{ID: uuid.New(), OriginalLink: original, ShortLink: short}, nil
		},
	})

	resp := do(t, app, "GET", "/"+short, nil)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, original, decode(t, resp)["original_url"])
}

func TestGetHandler_NotFound(t *testing.T) {
	app := setupApp(&mockService{
		getFunc: func(_ context.Context, _ string) (domain.Link, error) {
			return domain.Link{}, model.ErrNotFound
		},
	})

	resp := do(t, app, "GET", "/nonexistent", nil)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	assert.Equal(t, "short link not found", decode(t, resp)["error"])
}

func TestGetHandler_ServiceError(t *testing.T) {
	app := setupApp(&mockService{
		getFunc: func(_ context.Context, _ string) (domain.Link, error) {
			return domain.Link{}, errors.New("db error")
		},
	})

	resp := do(t, app, "GET", "/abc1234567", nil)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, "failed to retrieve original URL", decode(t, resp)["error"])
}
