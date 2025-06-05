package interfaces

import (
	"go-lti/internal/domain/dto"

	"github.com/gofiber/fiber/v2"
)

type CanvasService interface {
	Oauth2Login(c *fiber.Ctx) (string, error)
	Oauth2Redirect(c *fiber.Ctx, request *dto.Oauth2RedirectRequest) (*dto.Oauth2ExchangeResponse, error)
	Oauth2Refresh(c *fiber.Ctx, refreshToken string) (string, error)
	GetUserInfo(c *fiber.Ctx, accessToken string) (any, error)
}
