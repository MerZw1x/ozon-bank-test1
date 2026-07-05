package handler

import (
	"backend/src/internal/dto"
	"backend/src/internal/service"
	"net/url"

	"github.com/gofiber/fiber/v2"
)

type LinksHandler struct {
	service service.ILinkService
}

func NewLinkHandler(service service.ILinkService) *LinksHandler {
	return &LinksHandler{service: service}
}

func (h *LinksHandler) Shorten(c *fiber.Ctx) error {
	var req dto.SaveLinkRequest
	var res dto.SaveLinkResponse

	if err := c.BodyParser(&req); err != nil {
		res.Error = "invalid request body"
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	if req.URL == "" {
		res.Error = "url is required"
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	parsedURL, err := url.ParseRequestURI(req.URL)
	if err != nil || parsedURL.Scheme == "" {
		res.Error = "invalid URL format"
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	link, err := h.service.SaveLink(c.Context(), req.URL)
	if err != nil {
		res.Error = "failed to shorten URL: " + err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(res)
	}

	status := fiber.StatusOK

	res.ShortLink = link.ShortLink
	return c.Status(status).JSON(res)
}

func (h *LinksHandler) Redirect(c *fiber.Ctx) error {
	var res dto.RedirectResponse

	shortLink := c.Params("shortLink")
	if shortLink == "" {
		res.Error = "short link is required"
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	link, err := h.service.GetLink(c.Context(), shortLink)
	if err != nil {
		res.Error = "failed to retrieve original URL"
		return c.Status(fiber.StatusInternalServerError).JSON(res)
	}

	if link == nil {
		res.Error = "short link not found"
		return c.Status(fiber.StatusNotFound).JSON(res)
	}

	res.OriginalLink = link.OriginalLink
	return c.JSON(res)
}
