package interfaces

import (
	"go-lti/internal/domain/dto"

	"github.com/gofiber/fiber/v2"
)

type LtiService interface {
	// GetJwks returns the JSON Web Key Set for JWT validation
	GetJwks(c *fiber.Ctx) (*dto.JwksResponse, error)
	
	// LtiLogin handles the LTI login process and returns authorization URL
	LtiLogin(c *fiber.Ctx, request *dto.LtiLoginRequest) (string, error)
	
	LtiLaunch(c *fiber.Ctx, request *dto.LtiLaunchRequest) (*dto.LtiLaunchResponse, error)
	
	GetUserEmail(claims *dto.LtiJwtTokenClaims) string
	
	GetUserData(claims *dto.LtiJwtTokenClaims) *dto.LtiPrivacyClaims

	GenerateJwtToken(claims *dto.LtiPrivacyClaims) (string, error)

}