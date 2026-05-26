package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	Role         string `json:"role"`
	TokenVersion int    `json:"token_version"`
	TokenType    int    `json:"token_type"` // 1 - Access | 2 - Refresh
}

type AccessToken struct {
	Role         string `json:"role"`
	TokenVersion int    `json:"token_version"`
	jwt.RegisteredClaims
}

type RefreshTokenData struct {
	Id        string
	UserId    string
	Version   int
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
	Scope        []string  `json:"scope,omitempty"`
}

type UnSuccessfulResponse struct {
	Error            string `json:"error"` // invalid_request, invalid_client, invalid_grant, invalid_scope, unauthorised_client, unsupported_grant
	ErrorDescription string `json:"error_description,omitempty"`
	ErrorUri         string `json:"error_uri,omitempty"`
}

type JwtService struct {
}

func (js *JwtService) GenerateAccessToken(ctx context.Context, userid string) (string, error) {

}

func (js *JwtService) GenerateRefreshToken(ctx context.Context, userid string) (string, error) {

}

func (js *JwtService) GenerateTokenPair(ctx context.Context, userid string) (TokenPair, error) {

}

func (js *JwtService) ValidateAccessToken(ctx context.Context, tokenString string) (AccessToken, error) {

}

func (js *JwtService) ValidateRefreshToken(ctx context.Context, tokenString string) bool {

}

func (js *JwtService) RefreshTokens(ctx context.Context, refreshTokenString, role string) (TokenPair, error) {

}

func (js *JwtService) RevokeAccessToken(ctx context.Context) error {

}

func (js *JwtService) RevokeAllTokens(ctx context.Context) error {

}
