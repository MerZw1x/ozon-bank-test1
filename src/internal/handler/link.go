package handler

import (
	"backend/src/internal/domain"
	"backend/src/internal/model"
	"context"
	"errors"
	"log/slog"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gofiber/fiber/v2"
)

type LinksService interface {
	Shorten(ctx context.Context, originalLink string) (domain.Link, error)
	GetOriginal(ctx context.Context, shortLink string) (domain.Link, error)
}

type LinksHandler struct {
	service LinksService
	log     *slog.Logger
}

func NewLinksHandler(service LinksService, log *slog.Logger) *LinksHandler {
	return &LinksHandler{service: service, log: log}
}

func (h *LinksHandler) Register(app *fiber.App) {
	app.Post("/shorten", h.Shorten)
	app.Get("/:shortLink", h.Get)
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ShortLink string `json:"short_link"`
}

type getResponse struct {
	OriginalLink string `json:"original_url"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func validateShortenReq(req shortenRequest) error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.URL, validation.Required, is.URL),
	)
}

func (h *LinksHandler) Shorten(c *fiber.Ctx) error {
	var req shortenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{Error: "invalid request body"})
	}

	if err := validateShortenReq(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{Error: err.Error()})
	}

	link, err := h.service.Shorten(c.Context(), req.URL)
	if err != nil {
		h.log.Error("shorten failed", "url", req.URL, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse{Error: "failed to shorten URL"})
	}

	return c.Status(fiber.StatusOK).JSON(shortenResponse{ShortLink: link.ShortLink})
}

func (h *LinksHandler) Get(c *fiber.Ctx) error {
	shortLink := c.Params("shortLink")

	link, err := h.service.GetOriginal(c.Context(), shortLink)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(errorResponse{Error: "short link not found"})
		}
		h.log.Error("get original failed", "short_link", shortLink, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse{Error: "failed to retrieve original URL"})
	}

	return c.JSON(getResponse{OriginalLink: link.OriginalLink})
}
