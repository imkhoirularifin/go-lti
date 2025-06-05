package lti

import (
	"go-lti/internal/domain/dto"
	"go-lti/internal/domain/interfaces"

	"github.com/gofiber/fiber/v2"
)

type httpHandler struct {
	ltiService interfaces.LtiService
}

func NewHttpHandler(r fiber.Router, ltiService interfaces.LtiService) {
	handler := &httpHandler{
		ltiService: ltiService,
	}

	r.Post("/login", handler.ltiLogin)
	r.Post("/launch", handler.ltiLaunch)
	r.Get("/jwks", handler.jwks)
}

func (h *httpHandler) ltiLogin(c *fiber.Ctx) error {
	request := new(dto.LtiLoginRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ResponseDto{
			Message:     "Invalid request body",
			ErrorDetail: err.Error(),
		})
	}

	authURL, err := h.ltiService.LtiLogin(c, request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ResponseDto{
			Message:     "Failed to generate auth URL",
			ErrorDetail: err.Error(),
		})
	}

	return c.Redirect(authURL, fiber.StatusFound)
}

func (h *httpHandler) ltiLaunch(c *fiber.Ctx) error {
	request := new(dto.LtiLaunchRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ResponseDto{
			Message:     "Invalid request body",
			ErrorDetail: err.Error(),
		})
	}

	claims, err := h.ltiService.LtiLaunch(c, request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ResponseDto{
			Message:     "Failed to launch LTI",
			ErrorDetail: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.ResponseDto{
		Message: "LTI launch",
		Data:    claims,
	})
}

func (h *httpHandler) jwks(c *fiber.Ctx) error {
	jwks, err := h.ltiService.GetJwks(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ResponseDto{
			Message:     "Failed to get JWKS",
			ErrorDetail: err.Error(),
		})
	}

	return c.JSON(jwks)
}
