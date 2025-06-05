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
	r.Get("/access_token", handler.requestAccessToken)
}

func (h *httpHandler) jwks(c *fiber.Ctx) error {
	jwks, err := h.ltiService.GetJwks(c)
	if err != nil {
		return err
	}

	return c.JSON(jwks)
}

func (h *httpHandler) ltiLogin(c *fiber.Ctx) error {
	request := new(dto.LtiLoginRequest)
	if err := c.BodyParser(request); err != nil {
		return err
	}

	authURL, err := h.ltiService.LtiLogin(c, request)
	if err != nil {
		return err
	}

	return c.Redirect(authURL, fiber.StatusTemporaryRedirect)
}

func (h *httpHandler) ltiLaunch(c *fiber.Ctx) error {
	request := new(dto.LtiLaunchRequest)
	if err := c.BodyParser(request); err != nil {
		return err
	}

	claims, err := h.ltiService.LtiLaunch(c, request)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.ResponseDto{
		Message: "LTI launch",
		Data:    claims,
	})
}

func (h *httpHandler) requestAccessToken(c *fiber.Ctx) error {
	accessToken, err := h.ltiService.RequestAccessToken(c)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.ResponseDto{
		Message: "LTI access token",
		Data:    accessToken,
	})
}
