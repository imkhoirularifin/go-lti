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
	"go-lti/lib/httpclient"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type service struct {
	cfg        config.AppConfig
	httpClient httpclient.HttpClient
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
	key.Set(jwk.KeyIDKey, s.cfg.LtiConfig.JwkKid)
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
	fmt.Printf("request.IdToken: %v\n", request.IdToken)

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

// RequestAccessToken : Used to request LTI access token from Canvas
func (s *service) RequestAccessToken(c *fiber.Ctx) (any, error) {
	canvasDomain := s.cfg.CanvasConfig.Domain
	grantType := "client_credentials"
	clientAssertionType := "urn:ietf:params:oauth:client-assertion-type:jwt-bearer"
	clientAssertion, err := s.generateJWT()
	if err != nil {
		return nil, err
	}
	scope := "https://purl.imsglobal.org/spec/lti/scope/noticehandlers"

	url := fmt.Sprintf("https://%s/login/oauth2/token", canvasDomain)

	body := map[string]string{
		"grant_type":            grantType,
		"client_assertion_type": clientAssertionType,
		"client_assertion":      clientAssertion,
		"scope":                 scope,
	}

	var accessTokenResponse interface{}
	err = s.httpClient.Call(c.Context(), http.MethodPost, url, map[string]string{
		fiber.HeaderContentType: fiber.MIMEApplicationJSON,
		fiber.HeaderAccept:      fiber.MIMEApplicationJSON,
	}, body, &accessTokenResponse)
	if err != nil {
		return nil, err
	}

	return accessTokenResponse, nil
}

// validateJWT : Private method to validate JWT
func (s *service) validateJWT(idToken string) (*dto.LtiJwtTokenClaims, error) {
	jwksUrl := fmt.Sprintf("https://%s/api/lti/security/jwks", s.cfg.CanvasConfig.Domain)

	// Get JWKS with http client
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(jwksUrl)
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

// generateJWT : Private method to generate JWT for LTI access token request
func (s *service) generateJWT() (string, error) {
	privateKeyData, err := os.ReadFile(s.cfg.KeyConfig.PrivateKeyPath)
	if err != nil {
		return "", fmt.Errorf("failed to read private key: %w", err)
	}

	// Decode PEM
	block, _ := pem.Decode(privateKeyData)
	if block == nil {
		return "", fmt.Errorf("failed to decode PEM block")
	}

	// Parse private key
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	// Create JWT
	token := jwt.New()
	token.Set(jwt.IssuerKey, s.cfg.LtiConfig.Issuer)
	token.Set(jwt.SubjectKey, s.cfg.LtiConfig.ClientId)
	token.Set(jwt.AudienceKey, fmt.Sprintf("https://%s/login/oauth2/token", s.cfg.CanvasConfig.Domain))
	token.Set(jwt.IssuedAtKey, time.Now().Unix())
	token.Set(jwt.ExpirationKey, time.Now().Add(10*time.Minute).Unix())
	token.Set(jwt.JwtIDKey, uuid.New().String())

	// Create a key from the private key
	key, err := jwk.FromRaw(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to create key: %w", err)
	}
	key.Set(jwk.KeyIDKey, s.cfg.LtiConfig.JwkKid)
	key.Set(jwk.AlgorithmKey, jwa.RS256)
	key.Set(jwk.KeyUsageKey, "sig")

	// Sign the token with the key
	signedToken, err := jwt.Sign(token, jwt.WithKey(jwa.RS256, key))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return string(signedToken), nil
}

func NewService(
	cfg config.AppConfig,
	httpClient httpclient.HttpClient,
) interfaces.LtiService {
	return &service{
		cfg:        cfg,
		httpClient: httpClient,
		nonceCache: make(map[string]string),
	}
}
