package auth

//go:generate mockgen -destination=token_blacklist_gomock.go -package=auth -source=token_blacklist.go TokenBlacklistInterface -typed

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// TokenBlacklist is intentionally minimal to ease mocking in services
// and avoid leakage of implementation-specific methods
// Only the methods used by business services are exposed here.
type TokenBlacklist interface {
	// TestConnection verifies the Redis connection is healthy
	TestConnection() error

	BlacklistToken(ctx context.Context, tokenID string, expiresAt time.Time) error
	IsTokenBlacklisted(ctx context.Context, tokenID string) (bool, error)
}

type tokenBlacklist struct {
	client *redis.Client
}

type TokenBlacklistConfig struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

type BlacklistStats struct {
	BlacklistedTokens int
	RedisConnected    bool
}

func NewTokenBlacklist(config TokenBlacklistConfig) TokenBlacklist {
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	return &tokenBlacklist{
		client: client,
	}
}

func (tb *tokenBlacklist) TestConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := tb.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	return nil
}

// BlacklistToken adds a token to the blacklist with TTL
func (tb *tokenBlacklist) BlacklistToken(ctx context.Context, tokenID string, expiresAt time.Time) error {
	// Calculate TTL - how long until the token would naturally expire
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		// Token already expired, no need to blacklist
		return nil
	}

	key := fmt.Sprintf("blacklist:%s", tokenID)
	err := tb.client.Set(ctx, key, "1", ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}

	return nil
}

func (tb *tokenBlacklist) IsTokenBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	key := fmt.Sprintf("blacklist:%s", tokenID)
	result := tb.client.Exists(ctx, key)
	if result.Err() != nil {
		return false, fmt.Errorf("failed to check token blacklist: %w", result.Err())
	}

	return result.Val() > 0, nil
}

func (tb *tokenBlacklist) Close() error {
	return tb.client.Close()
}

func (tb *tokenBlacklist) GetStats(ctx context.Context) (BlacklistStats, error) {
	keys, err := tb.client.Keys(ctx, "blacklist:*").Result()
	if err != nil {
		return BlacklistStats{}, fmt.Errorf("failed to get blacklist stats: %w", err)
	}

	stats := BlacklistStats{
		BlacklistedTokens: len(keys),
		RedisConnected:    true,
	}

	return stats, nil
}
