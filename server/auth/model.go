package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

type NovaClaims struct {
	Role         string `json:"role"`
	TokenVersion int    `json:"token_version"`
	TokenType    int    `json:"token_type"` // 1 - Access | 2 - Refresh
	jwt.RegisteredClaims
}

type RefreshTokenData struct {
	Id        string
	UserId    string
	Version   int
	CreatedAt time.Time
	ExpiresAt time.Time
	IsRevoked bool
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SuccessfulResponse struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    time.Time `json:"expires_in"`
	RefreshToken string    `json:"refresh_token"`
	Scope        string    `json:"scope,omitempty"`
}

type LoginRequest struct {
	Provider string `json:"provider"`
	Code     string `json:"code"`
}

type LoginResponse struct {
	Status string              `json:"status"`
	Reason string              `json:"reason,omitempty"`
	Meta   map[string]string   `json:"meta,omitempty"`
	Tokens *SuccessfulResponse `json:"tokens,omitempty"`
}

type TokenRefreshRequest struct {
	GrantType    string `json:"grant_type" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type FailedRefreshResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
}

type TokenRefreshResponse struct {
	Successful *SuccessfulResponse
	Failed     *FailedRefreshResponse
}

type TokenResponse struct {
	oauth2.Token
	IdToken string `json:"id_token"`
	Scopes  string `json:"scope"`
}

type UnifiedUserProfile struct {
	Provider    string `json:"provider"`
	UserId      string `json:"userid"`
	DisplayName string `json:"name"`
	Email       string `json:"email"`
	AvatarUrl   string `json:"avatar_url"`
}

type GoogleProfile struct {
	Sub           string `json:"sub"` // Google's unique, permanent user ID
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"` // Direct URL to avatar image
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Locale        string `json:"locale"`
}

type FacebookProfile struct {
	ID        string `json:"id"` // App-scoped user ID
	Name      string `json:"name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Picture   struct {
		Data struct {
			URL          string `json:"url"` // Nested profile image URL
			IsSilhouette bool   `json:"is_silhouette"`
			Width        int    `json:"width"`
			Height       int    `json:"height"`
		} `json:"data"`
	} `json:"picture"`
}

type TwitterProfile struct {
	Data struct {
		ID              string `json:"id"`       // Permanent numerical string ID
		Name            string `json:"name"`     // Display name
		Username        string `json:"username"` // Twitter handle (e.g., @john)
		ProfileImageURL string `json:"profile_image_url"`
	} `json:"data"`
	// Note: Email often returns parallel to "data" if granted permissions
	Email string `json:"email,omitempty"`
}

type DiscordProfile struct {
	ID            string `json:"id"`            // Unique snowflake string ID
	Username      string `json:"username"`      // Clean login username
	GlobalName    string `json:"global_name"`   // Display name seen in chat
	Avatar        string `json:"avatar"`        // Hash string used to construct avatar URL
	Discriminator string `json:"discriminator"` // "0" for new accounts or legacy 4-digit tag
	Email         string `json:"email"`
	Verified      bool   `json:"verified"`
	Banner        string `json:"banner"`       // Banner hash string
	AccentColor   int    `json:"accent_color"` // Integer representation of color hex
}

type AppleClaims struct {
	Issuer         string `json:"iss"` // Always "https://appleid.apple.com"
	Subject        string `json:"sub"` // User's team-scoped unique identifier
	Audience       string `json:"aud"` // Your client_id / App ID
	ExpiresAt      int64  `json:"exp"`
	IssuedAt       int64  `json:"iat"`
	Email          string `json:"email"`
	EmailVerified  string `json:"email_verified"`   // Crucial: Apple returns this as a string ("true"/"false") or boolean depending on platform version
	IsPrivateEmail string `json:"is_private_email"` // Tells you if they used Apple Hide My Email
}
