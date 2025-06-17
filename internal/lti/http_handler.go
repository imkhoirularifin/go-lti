package lti

import (
	"fmt"
	"go-lti/internal/domain/dto"
	"go-lti/internal/domain/interfaces"
	"log"

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

	log.Printf("LTI Login Request - Issuer: %s, Client ID: %s, Login Hint: %s",
		request.Iss, request.ClientId, request.LoginHint)

	authURL, err := h.ltiService.LtiLogin(c, request)
	if err != nil {
		log.Printf("LTI Login failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ResponseDto{
			Message:     "Failed to generate auth URL",
			ErrorDetail: err.Error(),
		})
	}

	log.Printf("LTI Login successful, redirecting to: %s", authURL)
	return c.Redirect(authURL, fiber.StatusFound)
}

func (h *httpHandler) ltiLaunch(c *fiber.Ctx) error {
	request := new(dto.LtiLaunchRequest)
	if err := c.BodyParser(request); err != nil {
		log.Printf("Failed to parse LTI Launch request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.ResponseDto{
			Message:     "Invalid request body",
			ErrorDetail: err.Error(),
		})
	}

	log.Printf("LTI Launch Request - State: %s, Has ID Token: %t", request.State, request.IdToken != "")

	if request.Error != "" {
		log.Printf("LTI Launch error from Canvas: %s - %s", request.Error, request.ErrorDescription)
		return c.Status(fiber.StatusBadRequest).JSON(dto.ResponseDto{
			Message:     "LTI Launch error from Canvas",
			ErrorDetail: fmt.Sprintf("%s: %s", request.Error, request.ErrorDescription),
		})
	}

	launchResponse, err := h.ltiService.LtiLaunch(c, request)
	if err != nil {
		log.Printf("LTI Launch failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ResponseDto{
			Message:     "Failed to launch LTI",
			ErrorDetail: err.Error(),
		})
	}

	if launchResponse.Success && launchResponse.UserData != nil {
		log.Printf("LTI Launch successful for user: %s (Email: %s, Canvas ID: %s)",
			launchResponse.UserData.Name,
			launchResponse.UserData.Email,
			launchResponse.CanvasUserID,
		)

		token, err := h.ltiService.GenerateJwtToken(launchResponse.UserData)
		if err != nil {
			log.Printf("Failed to generate JWT token: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(dto.ResponseDto{
				Message:     "Failed to generate token",
				ErrorDetail: err.Error(),
			})
		}

		
		// redirectURL := fmt.Sprintf("http://localhost:3000/career-planner?token=%s", token)
		redirectURL := fmt.Sprintf("https://staging.app.ejourney.id/career-planner?token=%s", token)
		return c.Redirect(redirectURL, fiber.StatusTemporaryRedirect)
	}

	// Default response jika tidak memenuhi kondisi sukses
	return c.Status(fiber.StatusOK).JSON(dto.ResponseDto{
		Message: "LTI launch did not return user data",
		Data:    launchResponse,
	})
}

func (h *httpHandler) jwks(c *fiber.Ctx) error {
	jwks, err := h.ltiService.GetJwks(c)
	if err != nil {
		log.Printf("Failed to get JWKS: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ResponseDto{
			Message:     "Failed to get JWKS",
			ErrorDetail: err.Error(),
		})
	}

	return c.JSON(jwks)
}
