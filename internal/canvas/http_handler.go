package canvas

import (
	"go-lti/internal/domain/dto"
	"go-lti/internal/domain/interfaces"

	"github.com/gofiber/fiber/v2"
)

type httpHandler struct {
	canvasService interfaces.CanvasService
}

func NewHttpHandler(r fiber.Router, canvasService interfaces.CanvasService) {
	handler := &httpHandler{
		canvasService: canvasService,
	}

	r.Get("/oauth2/login", handler.Oauth2Login)
	r.Get("/oauth2/redirect", handler.Oauth2Redirect)
}

func (h *httpHandler) Oauth2Login(c *fiber.Ctx) error {
	loginUrl, err := h.canvasService.Oauth2Login(c)
	if err != nil {
		return err
	}

	return c.Redirect(loginUrl, fiber.StatusTemporaryRedirect)
}

func (h *httpHandler) Oauth2Redirect(c *fiber.Ctx) error {
	req := new(dto.Oauth2RedirectRequest)
	if err := c.QueryParser(req); err != nil {
		return err
	}

	exchangeResponse, err := h.canvasService.Oauth2Redirect(c, req)
	if err != nil {
		return err
	}

	userInfo, err := h.canvasService.GetUserInfo(c, exchangeResponse.AccessToken)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.ResponseDto{
		Message: "Successfully exchanged code for access token",
		Data:    userInfo,
	})
}
