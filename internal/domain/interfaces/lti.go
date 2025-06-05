package interfaces

import (
	"go-lti/internal/domain/dto"

	"github.com/gofiber/fiber/v2"
)

type LtiService interface {
	GetJwks(c *fiber.Ctx) (*dto.JwksResponse, error)
	LtiLogin(c *fiber.Ctx, request *dto.LtiLoginRequest) (string, error)
	LtiLaunch(c *fiber.Ctx, request *dto.LtiLaunchRequest) (*dto.LtiJwtTokenClaims, error)
}
