package auth

import (
	"context"
	"time"
)

type TokenStore interface {
	// StoreRefreshToken saves a refresh token with its metadata
	StoreRefreshToken(ctx context.Context, tokenID, userID string, expiresAt time.Time) error
	// GetRefreshToken retrieves refresh token metadata
	GetRefreshToken(ctx context.Context, tokenID string) (*RefreshTokenData, error)
	// RevokeRefreshToken marks a refresh token as revoked
	RevokeRefreshToken(ctx context.Context, tokenID string) error
	// RevokeAllUserTokens revokes all tokens for a specific user
	RevokeAllUserTokens(ctx context.Context, userID string) error
	// IsTokenBlacklisted checks if an access token has been blacklisted
	IsTokenBlacklisted(ctx context.Context, tokenID string) (bool, error)
	// BlacklistToken adds an access token to the blacklist
	BlacklistToken(ctx context.Context, tokenID string, expiresAt time.Time) error
	// GetUserTokenVersion returns the current token version for a user
	GetUserTokenVersion(ctx context.Context, userID string) (int, error)
	// IncrementUserTokenVersion increments and returns the new token version
	IncrementUserTokenVersion(ctx context.Context, userID string) (int64, error)
}
