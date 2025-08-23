package auth

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenBlacklist_Integration(t *testing.T) {
	t.Parallel()

	config := TokenBlacklistConfig{
		RedisAddr: "localhost:6379",
		RedisDB:   1, // Use DB 1 for tests
	}

	blacklist := NewTokenBlacklist(config)
	// Check connection and skip if Redis not available
	if err := blacklist.TestConnection(); err != nil {
		t.Skipf("Redis not available, skipping integration tests: %v", err)
		return
	}
	// Access implementation-only methods for cleanup and stats
	impl, ok := blacklist.(interface {
		Close() error
		GetStats(ctx context.Context) (BlacklistStats, error)
	})
	require.True(t, ok)
	defer impl.Close()

	ctx := context.Background()

	t.Run("it succeeds when blacklisting a token", func(t *testing.T) {
		tokenID := "test-token-123"
		expiresAt := time.Now().Add(time.Hour)

		err := blacklist.BlacklistToken(ctx, tokenID, expiresAt)
		assert.NoError(t, err)

		isBlacklisted, err := blacklist.IsTokenBlacklisted(ctx, tokenID)
		assert.NoError(t, err)
		assert.True(t, isBlacklisted)
	})

	t.Run("it succeeds when checking non-blacklisted token", func(t *testing.T) {
		tokenID := "non-blacklisted-token"

		isBlacklisted, err := blacklist.IsTokenBlacklisted(ctx, tokenID)
		assert.NoError(t, err)
		assert.False(t, isBlacklisted)
	})

	t.Run("it succeeds when token expires naturally", func(t *testing.T) {
		tokenID := "short-lived-token"
		expiresAt := time.Now().Add(100 * time.Millisecond)

		err := blacklist.BlacklistToken(ctx, tokenID, expiresAt)
		assert.NoError(t, err)

		isBlacklisted, err := blacklist.IsTokenBlacklisted(ctx, tokenID)
		assert.NoError(t, err)
		assert.True(t, isBlacklisted)

		time.Sleep(200 * time.Millisecond)

		isBlacklisted, err = blacklist.IsTokenBlacklisted(ctx, tokenID)
		assert.NoError(t, err)
		assert.False(t, isBlacklisted)
	})

	t.Run("it succeeds when blacklisting already expired token", func(t *testing.T) {
		tokenID := "already-expired-token"
		expiresAt := time.Now().Add(-time.Hour) // Already expired

		err := blacklist.BlacklistToken(ctx, tokenID, expiresAt)
		assert.NoError(t, err) // Should not error

		isBlacklisted, err := blacklist.IsTokenBlacklisted(ctx, tokenID)
		assert.NoError(t, err)
		assert.False(t, isBlacklisted)
	})

	t.Run("it succeeds when getting stats", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			tokenID := fmt.Sprintf("stats-test-token-%d", i)
			expiresAt := time.Now().Add(time.Hour)
			err := blacklist.BlacklistToken(ctx, tokenID, expiresAt)
			require.NoError(t, err)
		}

		impl := blacklist.(interface {
			GetStats(ctx context.Context) (BlacklistStats, error)
		})
		stats, err := impl.GetStats(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.True(t, stats.RedisConnected)
		assert.GreaterOrEqual(t, stats.BlacklistedTokens, 3)
	})
}

func TestTokenBlacklist_ConnectionFailure(t *testing.T) {
	t.Run("it fails when Redis is not available", func(t *testing.T) {
		config := TokenBlacklistConfig{
			RedisAddr: "localhost:9999", // Non-existent Redis
			RedisDB:   0,
		}

		bl := NewTokenBlacklist(config)
		err := bl.TestConnection()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to connect to Redis")
	})
}
