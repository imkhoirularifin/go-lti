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

// Privacy Claims
type LtiPrivacyClaims struct {
    UserID          string `json:"user_id,omitempty"`
    PersonSourcedID string `json:"person_sourcedid,omitempty"`
    Name            string `json:"name,omitempty"`
    GivenName       string `json:"given_name,omitempty"`
    FamilyName      string `json:"family_name,omitempty"`
    Email           string `json:"email,omitempty"`
    Picture         string `json:"picture,omitempty"`
    Locale          string `json:"locale,omitempty"`
    // Canvas specific fields
    CanvasUserID    string `json:"canvas_user_id,omitempty"`
    CanvasUserLogin string `json:"canvas_user_login_id,omitempty"`
    SisUserID       string `json:"sis_user_id,omitempty"`
}

// Canvas Custom Claims
type CanvasCustomClaims struct {
    CanvasUserID           string `json:"canvas_user_id,omitempty"`
    CanvasUserLoginID      string `json:"canvas_user_login_id,omitempty"`
    CanvasUserEmail        string `json:"canvas_user_email,omitempty"`
    CanvasUserName         string `json:"canvas_user_name,omitempty"`
    CanvasUserSisID        string `json:"canvas_user_sis_id,omitempty"`
    CanvasCourseID         string `json:"canvas_course_id,omitempty"`
    CanvasCourseName       string `json:"canvas_course_name,omitempty"`
    CanvasAssignmentID     string `json:"canvas_assignment_id,omitempty"`
    CanvasAssignmentTitle  string `json:"canvas_assignment_title,omitempty"`
    CanvasEnrollmentState  string `json:"canvas_enrollment_state,omitempty"`
    CanvasUserEnrollmentID string `json:"canvas_user_enrollment_id,omitempty"`
}

// Names and Role Provisioning Service (NRPS) Claims
type NrpsClaims struct {
    ContextMembershipsURL string   `json:"context_memberships_url"`
    ServiceVersions       []string `json:"service_versions"`
}

// Assignment and Grade Service (AGS) Claims
type AgsClaims struct {
    Scope    []string `json:"scope"`
    LineItem string   `json:"lineitems,omitempty"`
}

// Deep Linking Claims
type DeepLinkingClaims struct {
    DeepLinkingSettings struct {
        AcceptTypes                       []string `json:"accept_types"`
        AcceptMediaTypes                  string   `json:"accept_media_types"`
        AcceptPresentationDocumentTargets []string `json:"accept_presentation_document_targets"`
        AcceptMultiple                    bool     `json:"accept_multiple"`
        AutoCreate                        bool     `json:"auto_create"`
        Title                             string   `json:"title"`
        Text                              string   `json:"text"`
        Data                              string   `json:"data"`
    } `json:"https://purl.imsglobal.org/spec/lti-dl/claim/deep_linking_settings"`
}

// Modified LTI JWT Token Claims structure with explicit date fields as strings
type LtiJwtTokenClaims struct {
    // Standard JWT Claims - explicitly defined
    Issuer     string   `json:"iss"`
    Subject    string   `json:"sub"`
    Audience   []string `json:"aud,omitempty"`
    AudienceStr string  `json:"aud,omitempty"`
    Expiration string   `json:"exp"`
    NotBefore  string   `json:"nbf,omitempty"`
    IssuedAt   string   `json:"iat"`
    JWTID      string   `json:"jti,omitempty"`

    // Standard OpenID Connect Claims
    Nonce string `json:"nonce"`

    // Core LTI Claims
    MessageType  string `json:"https://purl.imsglobal.org/spec/lti/claim/message_type"`
    Version      string `json:"https://purl.imsglobal.org/spec/lti/claim/version"`
    ResourceLink struct {
        ID          string  `json:"id"`
        Description *string `json:"description"`
        Title       string  `json:"title"`
        URL         string  `json:"url,omitempty"`
    } `json:"https://purl.imsglobal.org/spec/lti/claim/resource_link"`

    // Platform and Context
    DeploymentID  string `json:"https://purl.imsglobal.org/spec/lti/claim/deployment_id"`
    TargetLinkURI string `json:"https://purl.imsglobal.org/spec/lti/claim/target_link_uri"`
    Context       struct {
        ID    string   `json:"id"`
        Title string   `json:"title"`
        Type  []string `json:"type"`
        Label string   `json:"label,omitempty"`
    } `json:"https://purl.imsglobal.org/spec/lti/claim/context"`

    ToolPlatform struct {
        GUID              string `json:"guid"`
        Name              string `json:"name"`
        Version           string `json:"version"`
        ProductFamilyCode string `json:"product_family_code"`
        ContactEmail      string `json:"contact_email,omitempty"`
        Description       string `json:"description,omitempty"`
        URL               string `json:"url,omitempty"`
    } `json:"https://purl.imsglobal.org/spec/lti/claim/tool_platform"`

    // Launch Presentation
    LaunchPresentation struct {
        DocumentTarget string `json:"document_target"`
        ReturnURL      string `json:"return_url"`
        Locale         string `json:"locale"`
        Height         int    `json:"height"`
        Width          int    `json:"width"`
    } `json:"https://purl.imsglobal.org/spec/lti/claim/launch_presentation"`

    // User Information
    Roles  []string `json:"https://purl.imsglobal.org/spec/lti/claim/roles"`
    Locale string   `json:"locale"`

    // PRIVACY CLAIMS
    ForUser *LtiPrivacyClaims `json:"https://purl.imsglobal.org/spec/lti/claim/for_user,omitempty"`

    // Legacy Support
    Lti11LegacyUserID string `json:"https://purl.imsglobal.org/spec/lti/claim/lti11_legacy_user_id"`
    Lti1p1            struct {
        UserID string `json:"user_id"`
    } `json:"https://purl.imsglobal.org/spec/lti/claim/lti1p1"`

    // Custom Claims (Canvas specific)
    Custom       interface{}         `json:"https://purl.imsglobal.org/spec/lti/claim/custom"`
    CanvasCustom *CanvasCustomClaims `json:"https://canvas.instructure.com/lti/claim/custom,omitempty"`

    // Canvas Specific
    Placement string `json:"https://www.instructure.com/placement"`

    // Platform Notification Service
    PlatformNotificationService *struct {
        ServiceVersions         []string `json:"service_versions"`
        PlatformNotificationURL string   `json:"platform_notification_service_url"`
        Scope                   []string `json:"scope"`
        NoticeTypesSupported    []string `json:"notice_types_supported"`
    } `json:"https://purl.imsglobal.org/spec/lti/claim/platformnotificationservice,omitempty"`

    // Names and Role Provisioning Service
    Nrps *NrpsClaims `json:"https://purl.imsglobal.org/spec/lti-nrps/claim/namesroleservice,omitempty"`

    // Assignment and Grade Service
    Ags *AgsClaims `json:"https://purl.imsglobal.org/spec/lti-ags/claim/endpoint,omitempty"`

    // Deep Linking
    DeepLinking *DeepLinkingClaims `json:",omitempty"`

    // Mentor relationship
    Mentor []string `json:"https://purl.imsglobal.org/spec/lti/claim/role_scope_mentor,omitempty"`
}

// Response DTOs
type LtiLaunchResponse struct {
    Success      bool                `json:"success"`
    UserData     *LtiPrivacyClaims   `json:"user_data,omitempty"`
    LtiUserID    string              `json:"lti_user_id"`
    CanvasUserID string              `json:"canvas_user_id,omitempty"`
    Roles        []string            `json:"roles,omitempty"`
    Context      *LtiContextInfo     `json:"context,omitempty"`
    CustomData   *CanvasCustomClaims `json:"custom_data,omitempty"`
    Error        string              `json:"error,omitempty"`
    ErrorDetails string              `json:"error_details,omitempty"`
}

type LtiContextInfo struct {
    ID    string   `json:"id"`
    Title string   `json:"title"`
    Type  []string `json:"type"`
    Label string   `json:"label,omitempty"`
}

// Configuration DTOs
type LtiToolConfiguration struct {
    Title             string                 `json:"title"`
    Description       string                 `json:"description"`
    OIDCInitiationURL string                 `json:"oidc_initiation_url"`
    TargetLinkURI     string                 `json:"target_link_uri"`
    Scopes            []string               `json:"scopes"`
    Extensions        []LtiExtension         `json:"extensions"`
    PublicJWK         map[string]interface{} `json:"public_jwk,omitempty"`
    PublicJWKURL      string                 `json:"public_jwk_url,omitempty"`
    CustomFields      map[string]string      `json:"custom_fields"`
    PrivacyLevel      string                 `json:"privacy_level"` // public, email_only, name_only, anonymous
}

type LtiExtension struct {
    Platform string            `json:"platform"`
    Settings LtiCanvasSettings `json:"settings"`
    Privacy  string            `json:"privacy_level"`
}

type LtiCanvasSettings struct {
    Platform            string         `json:"platform"`
    Placements          []LtiPlacement `json:"placements"`
    PrivacyLevel        string         `json:"privacy_level"`
    CourseNavigation    *LtiPlacement  `json:"course_navigation,omitempty"`
    AccountNavigation   *LtiPlacement  `json:"account_navigation,omitempty"`
    UserNavigation      *LtiPlacement  `json:"user_navigation,omitempty"`
    AssignmentSelection *LtiPlacement  `json:"assignment_selection,omitempty"`
    LinkSelection       *LtiPlacement  `json:"link_selection,omitempty"`
    PostGrades          bool           `json:"post_grades,omitempty"`
    UseCanvas           bool           `json:"use_1_3,omitempty"`
}

type LtiPlacement struct {
    Placement       string `json:"placement"`
    MessageType     string `json:"message_type"`
    TargetLinkURI   string `json:"target_link_uri"`
    Text            string `json:"text"`
    IconURL         string `json:"icon_url,omitempty"`
    CanvasIconClass string `json:"canvas_icon_class,omitempty"`
    WindowTarget    string `json:"window_target,omitempty"`
    SelectionHeight int    `json:"selection_height,omitempty"`
    SelectionWidth  int    `json:"selection_width,omitempty"`
    LaunchHeight    int    `json:"launch_height,omitempty"`
    LaunchWidth     int    `json:"launch_width,omitempty"`
    Enabled         bool   `json:"enabled"`
    Visibility      string `json:"visibility,omitempty"`
    Default         string `json:"default,omitempty"`   
}

// Utility DTOs
type LtiValidationResult struct {
    Valid     bool               `json:"valid"`
    Claims    *LtiJwtTokenClaims `json:"claims,omitempty"`
    UserData  *LtiPrivacyClaims  `json:"user_data,omitempty"`
    Error     string             `json:"error,omitempty"`
    ErrorCode string             `json:"error_code,omitempty"`
    Issuer    string             `json:"issuer,omitempty"`
    Audience  string             `json:"audience,omitempty"`
    Subject   string             `json:"subject,omitempty"`
    ExpiresAt int64              `json:"expires_at,omitempty"`
    IssuedAt  int64              `json:"issued_at,omitempty"`
}

// Canvas API Integration DTOs
type CanvasUserProfile struct {
    ID              int    `json:"id"`
    Name            string `json:"name"`
    ShortName       string `json:"short_name"`
    SortableName    string `json:"sortable_name"`
    Email           string `json:"email"`
    LoginID         string `json:"login_id"`
    SisUserID       string `json:"sis_user_id"`
    AvatarURL       string `json:"avatar_url"`
    Locale          string `json:"locale"`
    EffectiveLocale string `json:"effective_locale"`
    TimeZone        string `json:"time_zone"`
    Bio             string `json:"bio"`
}

type CanvasCourseInfo struct {
    ID               int    `json:"id"`
    Name             string `json:"name"`
    CourseCode       string `json:"course_code"`
    WorkflowState    string `json:"workflow_state"`
    AccountID        int    `json:"account_id"`
    StartAt          string `json:"start_at"`
    EndAt            string `json:"end_at"`
    EnrollmentTermID int    `json:"enrollment_term_id"`
    SisCourseID      string `json:"sis_course_id"`
}

// Error DTOs
type LtiError struct {
    Code        string `json:"code"`
    Message     string `json:"message"`
    Description string `json:"description,omitempty"`
    URI         string `json:"uri,omitempty"`
}

// Generic Response DTO
// type ResponseDto struct {
//     Message     string      `json:"message"`
//     ErrorDetail string      `json:"error_detail,omitempty"`
//     Data        interface{} `json:"data,omitempty"`
// }

// Constants untuk LTI Claims
const (
    // Privacy Levels
    PrivacyLevelPublic    = "public"
    PrivacyLevelEmailOnly = "email_only"
    PrivacyLevelNameOnly  = "name_only"
    PrivacyLevelAnonymous = "anonymous"

    // Message Types
    MessageTypeLtiResourceLinkRequest     = "LtiResourceLinkRequest"
    MessageTypeLtiDeepLinkingRequest      = "LtiDeepLinkingRequest"
    MessageTypeLtiSubmissionReviewRequest = "LtiSubmissionReviewRequest"

    // Placements
    PlacementCourseNavigation    = "course_navigation"
    PlacementAccountNavigation   = "account_navigation"
    PlacementUserNavigation      = "user_navigation"
    PlacementAssignmentSelection = "assignment_selection"
    PlacementLinkSelection       = "link_selection"
    PlacementRichEditor          = "editor_button"
    PlacementHomeworkSubmission  = "homework_submission"
    PlacementMigrationSelection  = "migration_selection"

    // Roles
    RoleInstructor        = "http://purl.imsglobal.org/vocab/lis/v2/membership#Instructor"
    RoleStudent           = "http://purl.imsglobal.org/vocab/lis/v2/membership#Learner"
    RoleTeachingAssistant = "http://purl.imsglobal.org/vocab/lis/v2/membership#TeachingAssistant"
    RoleAdmin             = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Administrator"
    RoleMentor            = "http://purl.imsglobal.org/vocab/lis/v2/membership#Mentor"
)

// Helper Methods
func (c *LtiJwtTokenClaims) HasPrivacyClaims() bool {
    return c.ForUser != nil
}

func (c *LtiJwtTokenClaims) GetUserEmail() string {
    if c.ForUser != nil && c.ForUser.Email != "" {
        return c.ForUser.Email
    }
    if c.CanvasCustom != nil && c.CanvasCustom.CanvasUserEmail != "" {
        return c.CanvasCustom.CanvasUserEmail
    }
    return ""
}

func (c *LtiJwtTokenClaims) GetUserName() string {
    if c.ForUser != nil && c.ForUser.Name != "" {
        return c.ForUser.Name
    }
    if c.CanvasCustom != nil && c.CanvasCustom.CanvasUserName != "" {
        return c.CanvasCustom.CanvasUserName
    }
    return ""
}

func (c *LtiJwtTokenClaims) GetCanvasUserID() string {
    if c.ForUser != nil && c.ForUser.CanvasUserID != "" {
        return c.ForUser.CanvasUserID
    }
    if c.CanvasCustom != nil && c.CanvasCustom.CanvasUserID != "" {
        return c.CanvasCustom.CanvasUserID
    }
    return ""
}

func (c *LtiJwtTokenClaims) IsInstructor() bool {
    for _, role := range c.Roles {
        if role == RoleInstructor || role == RoleTeachingAssistant {
            return true
        }
    }
    return false
}

func (c *LtiJwtTokenClaims) IsStudent() bool {
    for _, role := range c.Roles {
        if role == RoleStudent {
            return true
        }
    }
    return false
}

func (c *LtiJwtTokenClaims) IsAdmin() bool {
    for _, role := range c.Roles {
        if role == RoleAdmin {
            return true
        }
    }
    return false
}