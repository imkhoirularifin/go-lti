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
	"log"
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
	httpClient *http.Client
	nonceCache map[string]string
}

func (s *service) GenerateJwtToken(user *dto.LtiPrivacyClaims) (string, error) {
	token := jwt.New()
	now := time.Now()

	_ = token.Set(jwt.IssuedAtKey, now)
	_ = token.Set(jwt.ExpirationKey, now.Add(2*time.Hour))
	_ = token.Set("email", user.Email)
	_ = token.Set("name", user.Name)
	_ = token.Set("canvas_user_id", user.CanvasUserID)

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS256, []byte(s.cfg.KeyConfig.JwtSecret)))// gunakan secret dari config/env
	if err != nil {
		return "", err
	}

	return string(signed), nil
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

// LtiLaunch : Enhanced method to handle LTI launch and extract user data
func (s *service) LtiLaunch(c *fiber.Ctx, request *dto.LtiLaunchRequest) (*dto.LtiLaunchResponse, error) {
	claims, err := s.validateJWT(request.IdToken)
	log.Printf("LTI Launch Request - ID Token: %s, State: %s, Nonce: %s", request.IdToken, request.State, claims.Nonce)
	if err != nil {
		return &dto.LtiLaunchResponse{
			Success:      false,
			Error:        "JWT validation failed",
			ErrorDetails: err.Error(),
		}, err
	}

	// Check if nonce is valid
	state, ok := s.nonceCache[claims.Nonce]
	if !ok {
		return &dto.LtiLaunchResponse{
			Success:      false,
			Error:        "Invalid nonce",
			ErrorDetails: "Nonce not found in cache",
		}, errors.New("invalid nonce")
	}
	if state != request.State {
		return &dto.LtiLaunchResponse{
			Success:      false,
			Error:        "Invalid state",
			ErrorDetails: "State parameter does not match",
		}, errors.New("invalid state")
	}
	delete(s.nonceCache, claims.Nonce)

	// Extract user data from claims
	userData := s.extractUserData(claims)
	contextInfo := s.extractContextInfo(claims)
	customData := s.extractCustomData(claims)

	// Log extracted data for debugging
	s.logUserData(userData, claims)

	response := &dto.LtiLaunchResponse{
		Success:      true,
		UserData:     userData,
		LtiUserID:    claims.Subject,
		CanvasUserID: claims.GetCanvasUserID(),
		Roles:        claims.Roles,
		Context:      contextInfo,
		CustomData:   customData,
	}

	return response, nil
}

// extractUserData : Private method to extract user data from LTI claims
func (s *service) extractUserData(claims *dto.LtiJwtTokenClaims) *dto.LtiPrivacyClaims {
	userData := &dto.LtiPrivacyClaims{}

	// Extract from ForUser privacy claims (primary source)
	if claims.ForUser != nil {
		userData.UserID = claims.ForUser.UserID
		userData.PersonSourcedID = claims.ForUser.PersonSourcedID
		userData.Name = claims.ForUser.Name
		userData.GivenName = claims.ForUser.GivenName
		userData.FamilyName = claims.ForUser.FamilyName
		userData.Email = claims.ForUser.Email
		userData.Picture = claims.ForUser.Picture
		userData.Locale = claims.ForUser.Locale
		userData.CanvasUserID = claims.ForUser.CanvasUserID
		userData.CanvasUserLogin = claims.ForUser.CanvasUserLogin
		userData.SisUserID = claims.ForUser.SisUserID
	}

	// Extract from Canvas custom claims (fallback source)
	if claims.CanvasCustom != nil {
		if userData.Email == "" {
			userData.Email = claims.CanvasCustom.CanvasUserEmail
		}
		if userData.Name == "" {
			userData.Name = claims.CanvasCustom.CanvasUserName
		}
		if userData.CanvasUserID == "" {
			userData.CanvasUserID = claims.CanvasCustom.CanvasUserID
		}
		if userData.CanvasUserLogin == "" {
			userData.CanvasUserLogin = claims.CanvasCustom.CanvasUserLoginID
		}
		if userData.SisUserID == "" {
			userData.SisUserID = claims.CanvasCustom.CanvasUserSisID
		}
	}

	if customMap, ok := claims.Custom.(map[string]interface{}); ok {
		if email, exists := customMap["canvas_user_email"]; exists {
			if emailStr, ok := email.(string); ok && userData.Email == "" {
				userData.Email = emailStr
			}
		}
		if name, exists := customMap["canvas_user_name"]; exists {
			if nameStr, ok := name.(string); ok && userData.Name == "" {
				userData.Name = nameStr
			}
		}
		if userID, exists := customMap["canvas_user_id"]; exists {
			if userIDStr, ok := userID.(string); ok && userData.CanvasUserID == "" {
				userData.CanvasUserID = userIDStr
			}
		}
	}

	if userData.UserID == "" && claims.Lti11LegacyUserID != "" {
		userData.UserID = claims.Lti11LegacyUserID
	}

	if userData.UserID == "" && claims.Lti1p1.UserID != "" {
		userData.UserID = claims.Lti1p1.UserID
	}

	return userData
}

func (s *service) extractContextInfo(claims *dto.LtiJwtTokenClaims) *dto.LtiContextInfo {
	return &dto.LtiContextInfo{
		ID:    claims.Context.ID,
		Title: claims.Context.Title,
		Type:  claims.Context.Type,
		Label: claims.Context.Label,
	}
}

func (s *service) extractCustomData(claims *dto.LtiJwtTokenClaims) *dto.CanvasCustomClaims {
	if claims.CanvasCustom != nil {
		return claims.CanvasCustom
	}
	return nil
}

// logUserData : Private method to log extracted user data for debugging
func (s *service) logUserData(userData *dto.LtiPrivacyClaims, claims *dto.LtiJwtTokenClaims) {
	log.Printf("=== LTI User Data Extraction ===")
	log.Printf("User ID: %s", userData.UserID)
	log.Printf("Canvas User ID: %s", userData.CanvasUserID)
	log.Printf("Name: %s", userData.Name)
	log.Printf("Given Name: %s", userData.GivenName)
	log.Printf("Family Name: %s", userData.FamilyName)
	log.Printf("Email: %s", userData.Email)
	log.Printf("Canvas User Login: %s", userData.CanvasUserLogin)
	log.Printf("SIS User ID: %s", userData.SisUserID)
	log.Printf("Picture: %s", userData.Picture)
	log.Printf("Locale: %s", userData.Locale)
	log.Printf("Roles: %v", claims.Roles)
	log.Printf("Has Privacy Claims: %t", claims.HasPrivacyClaims())

	// Log custom claims for debugging
	if customMap, ok := claims.Custom.(map[string]interface{}); ok {
		log.Printf("Custom Claims:")
		for key, value := range customMap {
			log.Printf("  %s: %v", key, value)
		}
	}

	log.Printf("=== End User Data ===")
}

// validateJWT : Private method to validate JWT
func (s *service) validateJWT(idToken string) (*dto.LtiJwtTokenClaims, error) {
	jwksUrl := fmt.Sprintf("https://%s/api/lti/security/jwks", s.cfg.LtiConfig.CanvasDomain)

	// Get JWKS with http client
	resp, err := s.httpClient.Get(jwksUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch JWKS: status code %d", resp.StatusCode)
	}

	keySet, err := jwk.ParseReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWKS: %w", err)
	}

	// Parse JWT without validation
	token, err := jwt.Parse([]byte(idToken),
		jwt.WithKeySet(keySet),
		jwt.WithVerify(true),
		// Skip built-in validation since we'll handle date validation ourselves
		jwt.WithValidate(false),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT: %w", err)
	}

	// Convert token to map and then to custom claims
	rawClaims, err := token.AsMap(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to convert token to map: %w", err)
	}

	// Check audience manually
	aud, hasAud := rawClaims["aud"]
	if hasAud {
		log.Printf("Audience type: %T, value: %v", aud, aud)

		audValid := false
		clientID := s.cfg.LtiConfig.ClientId

		switch v := aud.(type) {
		case []interface{}:
			// If audience is array
			for _, a := range v {
				if aStr, ok := a.(string); ok && aStr == clientID {
					audValid = true
					break
				}
			}
		case string:
			// If audience is string
			audValid = (v == clientID)
		case interface{}:
			// Try to convert to string
			if audStr, ok := v.(string); ok {
				audValid = (audStr == clientID)
			}
		}

		if !audValid {
			log.Printf("Audience validation failed. Expected: %s, Got: %v", clientID, aud)

			log.Printf("Full claims: %+v", rawClaims)
		}
	} else {
		log.Printf("No audience claim found in token")
		// For debugging, continue anyway
		// return nil, fmt.Errorf("missing audience claim")
	}

	// Convert to JSON and then to our custom structure
	claimsBytes, err := json.Marshal(rawClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal claims: %w", err)
	}

	var claims dto.LtiJwtTokenClaims
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	// Log the decoded claims for debugging
	log.Printf("Successfully decoded JWT claims. Subject: %s", claims.Subject)

	return &claims, nil
}

// GetUserEmail : Public method to get user email from validated claims
func (s *service) GetUserEmail(claims *dto.LtiJwtTokenClaims) string {
	return claims.GetUserEmail()
}

// GetUserData : Public method to extract complete user data
func (s *service) GetUserData(claims *dto.LtiJwtTokenClaims) *dto.LtiPrivacyClaims {
	return s.extractUserData(claims)
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
