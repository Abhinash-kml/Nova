package auth

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisAuthStore struct {
	client *redis.Client
}

func NewRedisAuthStore(client *redis.Client) *RedisAuthStore {
	return &RedisAuthStore{client: client}
}

// StoreRefreshToken saves a refresh token with its metadata
func (s *RedisAuthStore) StoreRefreshToken(ctx context.Context, tokenID, userID string, expiresAt time.Time) error {

}

// GetRefreshToken retrieves refresh token metadata
func (s *RedisAuthStore) GetRefreshToken(ctx context.Context, tokenID string) (*RefreshTokenData, error) {

}

// RevokeRefreshToken marks a refresh token as revoked
func (s *RedisAuthStore) RevokeRefreshToken(ctx context.Context, tokenID string) error {

}

// RevokeAllUserTokens revokes all tokens for a specific user
func (s *RedisAuthStore) RevokeAllUserTokens(ctx context.Context, userID string) error {

}

// IsTokenBlacklisted checks if an access token has been blacklisted
func (s *RedisAuthStore) IsTokenBlacklisted(ctx context.Context, tokenID string) (bool, error) {

}

// BlacklistToken adds an access token to the blacklist
func (s *RedisAuthStore) BlacklistToken(ctx context.Context, tokenID string, expiresAt time.Time) error {

}

// GetUserTokenVersion returns the current token version for a user
func (s *RedisAuthStore) GetUserTokenVersion(ctx context.Context, userID string) (int, error) {

}

// IncrementUserTokenVersion increments and returns the new token version
func (s *RedisAuthStore) IncrementUserTokenVersion(ctx context.Context, userID string) (int64, error) {

}
