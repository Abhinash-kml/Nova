package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/abhinash-kml/nova/server/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Common errors for JWT operations
var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("token has expired")
	ErrInvalidClaims    = errors.New("invalid token claims")
	ErrTokenRevoked     = errors.New("token has been revoked")
	ErrInvalidTokenType = errors.New("invalid token type")
)

type CustomClaims struct {
	Role         string `json:"role"`
	TokenVersion int    `json:"token_version"`
	TokenType    int    `json:"token_type"` // 1 - Access | 2 - Refresh
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
	config *config.AuthTokenConfig
	store  TokenStore
}

func (js *JwtService) GenerateAccessToken(ctx context.Context, userid, role string) (string, error) {
	// Get the current token version for this user
	version, err := js.store.GetUserTokenVersion(ctx, userid)
	if err != nil {
		return "", fmt.Errorf("failed to get token version: %w", err)
	}

	now := time.Now()
	tokenID := uuid.New().String()

	claims := CustomClaims{
		Role:         role,
		TokenVersion: version,
		TokenType:    1,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID,
			Subject:   userid,
			Issuer:    js.config.AccessToken.Issuer,
			Audience:  jwt.ClaimStrings{js.config.AccessToken.Audience},
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(js.config.AccessToken.ExpiresIn))),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString(js.config.AccessToken.Secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (js *JwtService) GenerateRefreshToken(ctx context.Context, userid string) (string, error) {
	now := time.Now()
	tokenID := uuid.New().String()
	expiresAt := now.Add(time.Duration(js.config.RefreshToken.ExpiresIn))

	claims := CustomClaims{
		TokenType: 2,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID,
			Subject:   userid,
			Issuer:    js.config.RefreshToken.Issuer,
			Audience:  jwt.ClaimStrings{js.config.RefreshToken.Audience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedString, err := token.SignedString(js.config.RefreshToken.Secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	// Store refresh token in token store
	err = js.store.StoreRefreshToken(ctx, tokenID, userid, expiresAt)
	if err != nil {
		return "", fmt.Errorf("failed to store refresh token in tokenstore: %w", err)
	}

	return signedString, nil
}

func (js *JwtService) GenerateTokenPair(ctx context.Context, userid, role string) (*TokenPair, error) {
	access, err := js.GenerateAccessToken(ctx, userid, role)
	if err != nil {
		return nil, err
	}
	refresh, err := js.GenerateRefreshToken(ctx, userid)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (js *JwtService) ValidateAccessToken(ctx context.Context, tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(js.config.AccessToken.Secret), nil
	}, jwt.WithValidMethods([]string{"HS512"}),
		jwt.WithAllAudiences(js.config.AccessToken.Audience),
		jwt.WithIssuer(js.config.AccessToken.Issuer),
		jwt.WithExpirationRequired())
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidClaims
	}

	// Verify this is an access token, not a refresh token
	if claims.TokenType != 1 {
		return nil, ErrInvalidTokenType
	}

	// Check if the token has been blacklisted
	blacklisted, err := js.store.IsTokenBlacklisted(ctx, claims.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check token blacklist: %w", err)
	}
	if blacklisted {
		return nil, ErrTokenRevoked
	}

	// Verify token version matches current user version (for mass revocation)
	currentVersion, err := js.store.GetUserTokenVersion(ctx, claims.Subject)
	if err != nil {
		return nil, fmt.Errorf("failed to get user token version: %w", err)
	}
	if claims.TokenVersion < currentVersion {
		return nil, ErrTokenRevoked
	}

	return claims, nil
}

func (js *JwtService) ValidateRefreshToken(ctx context.Context, tokenString string) bool {

}

func (js *JwtService) RefreshTokens(ctx context.Context, refreshTokenString, role string) (TokenPair, error) {

}

func (js *JwtService) RevokeAccessToken(ctx context.Context) error {

}

func (js *JwtService) RevokeAllTokens(ctx context.Context) error {

}
