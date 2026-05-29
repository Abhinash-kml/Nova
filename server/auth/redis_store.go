package auth

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	refreshTokenPrefix = "refresh_token:"
	blacklistPrefix    = "blacklist:"
	tokenVersionPrefix = "token_version:"
	userTokenPrefix    = "user_tokens:"
)

type RedisAuthStore struct {
	client *redis.Client
}

func NewRedisAuthStore(client *redis.Client) *RedisAuthStore {
	return &RedisAuthStore{client: client}
}

// StoreRefreshToken saves a refresh token with its metadata
func (s *RedisAuthStore) StoreRefreshToken(ctx context.Context, tokenID, userID string, expiresAt time.Time) error {
	data := &RefreshTokenData{
		Id:        tokenID,
		UserId:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
		IsRevoked: false,
	}

	// Convert data to json
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Store with TTL matching the token expiration
	ttl := time.Until(expiresAt)
	err = s.client.Set(ctx, refreshTokenPrefix+tokenID, jsonData, ttl).Err()
	if err != nil {
		return err
	}

	// Add the token to user's token set for bulk revocation
	return s.client.SAdd(ctx, userTokenPrefix+userID, tokenID).Err()
}

// GetRefreshToken retrieves refresh token metadata
func (s *RedisAuthStore) GetRefreshToken(ctx context.Context, tokenID string) (*RefreshTokenData, error) {
	jsonData, err := s.client.Get(ctx, refreshTokenPrefix+tokenID).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var data RefreshTokenData
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// RevokeRefreshToken marks a refresh token as revoked
func (s *RedisAuthStore) RevokeRefreshToken(ctx context.Context, tokenID string) error {
	data, err := s.GetRefreshToken(ctx, tokenID)
	if err != nil || data == nil {
		return err
	}

	data.IsRevoked = true
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Keep the same TTL
	ttl := time.Until(data.ExpiresAt)
	return s.client.Set(ctx, refreshTokenPrefix+tokenID, jsonData, ttl).Err()
}

// RevokeAllUserTokens revokes all tokens for a specific user
func (s *RedisAuthStore) RevokeAllUserTokens(ctx context.Context, userID string) error {
	// Get all tokens for this user with tokenID
	tokenIDs, err := s.client.SMembers(ctx, userTokenPrefix+userID).Result()
	if err != nil {
		return err
	}

	// Revoke each token
	for _, token := range tokenIDs {
		err := s.RevokeRefreshToken(ctx, token)
		if err != nil {
			continue
		}
	}

	return nil
}

// IsTokenBlacklisted checks if an access token has been blacklisted
func (s *RedisAuthStore) IsTokenBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	exists, err := s.client.Exists(ctx, blacklistPrefix+tokenID).Result()
	return exists > 0, err
}

// BlacklistToken adds an access token to the blacklist
func (s *RedisAuthStore) BlacklistToken(ctx context.Context, tokenID string, expiresAt time.Time) error {
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		// Token has already expired, no need to blacklist it again
		return nil
	}
	return s.client.Set(ctx, blacklistPrefix+tokenID, "1", ttl).Err()
}

// GetUserTokenVersion returns the current token version for a user
func (s *RedisAuthStore) GetUserTokenVersion(ctx context.Context, userID string) (int, error) {
	version, err := s.client.Get(ctx, tokenVersionPrefix+userID).Int()
	if err == redis.Nil {
		return 0, nil // Default version is 0
	}
	return version, err
}

// IncrementUserTokenVersion increments and returns the new token version
func (s *RedisAuthStore) IncrementUserTokenVersion(ctx context.Context, userID string) (int64, error) {
	return s.client.Incr(ctx, tokenVersionPrefix+userID).Result()
}
