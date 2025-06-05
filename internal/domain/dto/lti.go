package dto

import (
	"github.com/lestrrat-go/jwx/v2/jwk"
)

type LtiLoginRequest struct {
	Iss               string `form:"iss"`
	LoginHint         string `form:"login_hint"`
	ClientId          string `form:"client_id"`
	LtiDeploymentId   string `form:"lti_deployment_id"`
	TargetLinkUri     string `form:"target_link_uri"`
	LtiMessageHint    string `form:"lti_message_hint"`
	CanvasEnvironment string `form:"canvas_environment"`
	CanvasRegion      string `form:"canvas_region"`
	LtiStorageTarget  string `form:"lti_storage_target"`
}

type LtiLaunchRequest struct {
	Utf8              string `form:"utf8"`
	AuthenticityToken string `form:"authenticity_token"`
	IdToken           string `form:"id_token"`
	State             string `form:"state"`
	LtiStorageTarget  string `form:"lti_storage_target"`
	Error             string `form:"error"`
	ErrorDescription  string `form:"error_description"`
}

type JwksResponse struct {
	Keys []jwk.Key `json:"keys"`
}

type LtiJwtTokenClaims struct {
	MessageType  string `json:"https://purl.imsglobal.org/spec/lti/claim/message_type"`
	Version      string `json:"https://purl.imsglobal.org/spec/lti/claim/version"`
	ResourceLink struct {
		ID          string  `json:"id"`
		Description *string `json:"description"`
		Title       string  `json:"title"`
	} `json:"https://purl.imsglobal.org/spec/lti/claim/resource_link"`
	Aud           []string `json:"aud"`
	Azp           string   `json:"azp"`
	DeploymentID  string   `json:"https://purl.imsglobal.org/spec/lti/claim/deployment_id"`
	Exp           string   `json:"exp"`
	Iat           string   `json:"iat"`
	Iss           string   `json:"iss"`
	Nonce         string   `json:"nonce"`
	Sub           string   `json:"sub"`
	TargetLinkURI string   `json:"https://purl.imsglobal.org/spec/lti/claim/target_link_uri"`
	Context       struct {
		ID    string   `json:"id"`
		Title string   `json:"title"`
		Type  []string `json:"type"`
	} `json:"https://purl.imsglobal.org/spec/lti/claim/context"`
	ToolPlatform struct {
		GUID              string `json:"guid"`
		Name              string `json:"name"`
		Version           string `json:"version"`
		ProductFamilyCode string `json:"product_family_code"`
	} `json:"https://purl.imsglobal.org/spec/lti/claim/tool_platform"`
	LaunchPresentation struct {
		DocumentTarget string `json:"document_target"`
		ReturnURL      string `json:"return_url"`
		Locale         string `json:"locale"`
		Height         int    `json:"height"`
		Width          int    `json:"width"`
	} `json:"https://purl.imsglobal.org/spec/lti/claim/launch_presentation"`
	PlatformNotificationService struct {
		ServiceVersions         []string `json:"service_versions"`
		PlatformNotificationURL string   `json:"platform_notification_service_url"`
		Scope                   []string `json:"scope"`
		NoticeTypesSupported    []string `json:"notice_types_supported"`
	} `json:"https://purl.imsglobal.org/spec/lti/claim/platformnotificationservice"`
	Locale   string   `json:"locale"`
	Roles    []string `json:"https://purl.imsglobal.org/spec/lti/claim/roles"`
	Custom   struct{} `json:"https://purl.imsglobal.org/spec/lti/claim/custom"`
	Endpoint struct {
		Scope     []string `json:"scope"`
		LineItems string   `json:"lineitems"`
	} `json:"https://purl.imsglobal.org/spec/lti-ags/claim/endpoint"`
	NamesRoleService struct {
		ContextMembershipsUrl string   `json:"context_memberships_url"`
		ServiceVersions       []string `json:"service_versions"`
	} `json:"https://purl.imsglobal.org/spec/lti-nrps/claim/namesroleservice"`
	Lti11LegacyUserID string `json:"https://purl.imsglobal.org/spec/lti/claim/lti11_legacy_user_id"`
	Lti1p1            struct {
		UserID string `json:"user_id"`
	} `json:"https://purl.imsglobal.org/spec/lti/claim/lti1p1"`
	Placement string `json:"https://www.instructure.com/placement"`
}
