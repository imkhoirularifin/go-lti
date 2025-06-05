package lti

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"go-lti/internal/domain/dto"
	"go-lti/internal/domain/interfaces"
	"go-lti/lib/config"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type service struct {
	cfg        config.AppConfig
	httpClient *http.Client
	nonceCache map[string]string
}

// GetJwks : Public method to return the JSON Web Key Set (JWKS) containing the public key used for JWT validation.
func (s *service) GetJwks(c *fiber.Ctx) (*dto.JwksResponse, error) {
	publicKey, err := os.ReadFile(s.cfg.KeyConfig.PublicKeyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(publicKey)
	rsaKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// Create JWK
	key, err := jwk.FromRaw(rsaKey)
	if err != nil {
		return nil, err
	}

	// Set JWK attributes
	key.Set(jwk.KeyIDKey, "my-lti-key")
	key.Set(jwk.AlgorithmKey, "RS256")
	key.Set(jwk.KeyUsageKey, "sig")

	return &dto.JwksResponse{
		Keys: []jwk.Key{key},
	}, nil
}

// LtiLogin : Public method to handle LTI login
func (s *service) LtiLogin(c *fiber.Ctx, request *dto.LtiLoginRequest) (string, error) {
	issuer := request.Iss
	scope := "openid"
	responseType := "id_token"
	clientId := request.ClientId
	redirectUri := s.cfg.LtiConfig.LaunchUrl
	loginHint := request.LoginHint
	ltiMessageHint := request.LtiMessageHint
	state := uuid.New().String()
	responseMode := "form_post"
	nonce := uuid.New().String()
	prompt := "none"

	// store nonce in cache
	s.nonceCache[nonce] = state

	authURL := fmt.Sprintf("%s/api/lti/authorize_redirect?scope=%s&response_type=%s&client_id=%s&redirect_uri=%s&login_hint=%s&lti_message_hint=%s&state=%s&response_mode=%s&nonce=%s&prompt=%s",
		issuer, scope, responseType, clientId, redirectUri, loginHint, ltiMessageHint, state, responseMode, nonce, prompt)

	return authURL, nil
}

// LtiLaunch : Public method to handle LTI launch
func (s *service) LtiLaunch(c *fiber.Ctx, request *dto.LtiLaunchRequest) (*dto.LtiJwtTokenClaims, error) {
	claims, err := s.validateJWT(request.IdToken)
	if err != nil {
		return nil, err
	}

	// Check if nonce is valid
	state, ok := s.nonceCache[claims.Nonce]
	if !ok {
		return nil, errors.New("invalid nonce")
	}
	if state != request.State {
		return nil, errors.New("invalid state")
	}
	delete(s.nonceCache, claims.Nonce)

	return claims, nil
}

// validateJWT : Private method to validate JWT
func (s *service) validateJWT(idToken string) (*dto.LtiJwtTokenClaims, error) {
	jwksUrl := fmt.Sprintf("https://%s/api/lti/security/jwks", s.cfg.LtiConfig.CanvasDomain)

	// Get JWKS with http client
	resp, err := s.httpClient.Get(jwksUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	keySet, err := jwk.ParseReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse and validate JWT
	token, err := jwt.Parse([]byte(idToken),
		jwt.WithKeySet(keySet),
		jwt.WithVerify(true),
		jwt.WithValidate(true),
		jwt.WithAudience(s.cfg.LtiConfig.ClientId),
	)
	if err != nil {
		return nil, err
	}

	// Convert token to LtiJwtTokenClaims
	rawClaims, err := token.AsMap(context.Background())
	if err != nil {
		return nil, err
	}

	claimsBytes, err := json.Marshal(rawClaims)
	if err != nil {
		return nil, err
	}

	var claims dto.LtiJwtTokenClaims
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return nil, err
	}

	return &claims, nil
}

func NewService(
	cfg config.AppConfig,
) interfaces.LtiService {
	httpClient := &http.Client{}

	return &service{
		cfg:        cfg,
		httpClient: httpClient,
		nonceCache: make(map[string]string),
	}
}
