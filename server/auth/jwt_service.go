package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	Role         string `json:"role"`
	TokenVersion int    `json:"token_version"`
	TokenType    int    `json:"token_type"` // 1 - Access | 2 - Refresh
}

type TokenPair struct {
	AccessToken  CustomClaims `json:"access_token"`
	RefreshToken CustomClaims `json:"refresh_token"`
}

type SuccessfulResponse struct {
	AccessToken  CustomClaims `json:"access_token"`
	TokenType    string       `json:"token_type"`
	ExpiresIn    time.Time    `json:"expires_in"`
	RefreshToken CustomClaims `json:"refresh_token"`
	Scope        []string     `json:"scope,omitempty"`
}

type UnSuccessfulResponse struct {
	Error            string `json:"error"` // invalid_request, invalid_client, invalid_grant, invalid_scope, unauthorised_client, unsupported_grant
	ErrorDescription string `json:"error_description,omitempty"`
	ErrorUri         string `json:"error_uri,omitempty"`
}

type JwtService struct {
}
