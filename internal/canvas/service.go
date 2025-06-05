package canvas

import (
	"fmt"
	"go-lti/internal/domain/dto"
	"go-lti/internal/domain/interfaces"
	"go-lti/lib/config"
	"go-lti/lib/httpclient"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type service struct {
	cfg        config.AppConfig
	httpClient httpclient.HttpClient
	stateCache map[string]string
}

// Oauth2Login : Redirect user to Canvas Oauth2 login page
func (s *service) Oauth2Login(c *fiber.Ctx) (string, error) {
	state := uuid.New().String()
	canvasDomain := s.cfg.CanvasConfig.Domain
	clientId := s.cfg.ApiKeyConfig.ClientId
	redirectUrl := s.cfg.ApiKeyConfig.RedirectUrl

	loginUrl := fmt.Sprintf("https://%s/login/oauth2/auth?client_id=%s&response_type=code&state=%s&redirect_uri=%s", canvasDomain, clientId, state, redirectUrl)

	s.stateCache[state] = state

	return loginUrl, nil
}

// Oauth2Redirect : Receive oauth2 callback from Canvas and exchange code for access token
func (s *service) Oauth2Redirect(c *fiber.Ctx, request *dto.Oauth2RedirectRequest) (*dto.Oauth2ExchangeResponse, error) {
	if request.Error != "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, request.ErrorDescription)
	}

	if request.State != s.stateCache[request.State] {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid state")
	}
	delete(s.stateCache, request.State)

	canvasDomain := s.cfg.CanvasConfig.Domain
	grantType := "authorization_code"
	clientId := s.cfg.ApiKeyConfig.ClientId
	clientSecret := s.cfg.ApiKeyConfig.Secret
	redirectUrl := s.cfg.ApiKeyConfig.RedirectUrl
	code := request.Code

	url := fmt.Sprintf("https://%s/login/oauth2/token?grant_type=%s&client_id=%s&client_secret=%s&code=%s&redirect_uri=%s", canvasDomain, grantType, clientId, clientSecret, code, redirectUrl)

	var exchangeResponse dto.Oauth2ExchangeResponse
	err := s.httpClient.Call(c.Context(), http.MethodPost, url, map[string]string{
		fiber.HeaderContentType: fiber.MIMEApplicationForm,
		fiber.HeaderAccept:      fiber.MIMEApplicationJSON,
	}, nil, &exchangeResponse)
	if err != nil {
		return nil, err
	}

	return &exchangeResponse, nil
}

// Oauth2Refresh : Used to get new access token using refresh token
func (s *service) Oauth2Refresh(c *fiber.Ctx, refreshToken string) (string, error) {
	panic("unimplemented")
}

// GetUserInfo : Used to get user info from Canvas
func (s *service) GetUserInfo(c *fiber.Ctx, accessToken string) (any, error) {
	url := fmt.Sprintf("https://%s/api/v1/users/self", s.cfg.CanvasConfig.Domain)

	var userInfo interface{}
	err := s.httpClient.Call(c.Context(), http.MethodGet, url, map[string]string{
		fiber.HeaderAuthorization: fmt.Sprintf("Bearer %s", accessToken),
	}, nil, &userInfo)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

func NewService(
	cfg config.AppConfig,
	httpClient httpclient.HttpClient,
) interfaces.CanvasService {
	return &service{
		cfg:        cfg,
		httpClient: httpClient,
		stateCache: make(map[string]string),
	}
}
