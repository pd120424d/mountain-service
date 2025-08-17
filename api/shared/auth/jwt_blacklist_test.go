package auth

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestValidateJWTWithBlacklist(t *testing.T) {
	t.Parallel()

	originalSecret := os.Getenv("JWT_SECRET")
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
	}()

	t.Run("it succeeds when token is valid and not blacklisted", func(t *testing.T) {
		token, err := GenerateJWT(1, "Employee")
		require.NoError(t, err)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		blacklist := NewMockTokenBlacklist(ctrl)
		blacklist.EXPECT().IsTokenBlacklisted(gomock.Any(), gomock.Any()).Return(false, nil).AnyTimes()
		blacklist.EXPECT().BlacklistToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

		claims, err := ValidateJWT(token, blacklist)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, uint(1), claims.ID)
		assert.Equal(t, "Employee", claims.Role)
	})

	t.Run("it succeeds when token is valid and blacklist is nil", func(t *testing.T) {
		token, err := GenerateJWT(1, "Employee")
		require.NoError(t, err)

		claims, err := ValidateJWT(token, nil)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, uint(1), claims.ID)
		assert.Equal(t, "Employee", claims.Role)
	})

	t.Run("it fails when token is blacklisted", func(t *testing.T) {
		token, err := GenerateJWT(1, "Employee")
		require.NoError(t, err)

		claims, err := ValidateJWT(token, nil)
		require.NoError(t, err)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		blacklist := NewMockTokenBlacklist(ctrl)
		// First call: mark blacklisted
		claimsID := claims.RegisteredClaims.ID
		blacklist.EXPECT().BlacklistToken(gomock.Any(), claimsID, gomock.Any()).Return(nil).AnyTimes()
		blacklist.EXPECT().IsTokenBlacklisted(gomock.Any(), claimsID).Return(true, nil).AnyTimes()

		validatedClaims, err := ValidateJWT(token, blacklist)
		assert.Error(t, err)
		assert.Nil(t, validatedClaims)
		assert.Contains(t, err.Error(), "token has been revoked")
	})

	t.Run("it fails when blacklist check returns error", func(t *testing.T) {
		token, err := GenerateJWT(1, "Employee")
		require.NoError(t, err)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		blacklist := NewMockTokenBlacklist(ctrl)
		blacklist.EXPECT().IsTokenBlacklisted(gomock.Any(), gomock.Any()).Return(false, fmt.Errorf("mock error")).AnyTimes()

		claims, err := ValidateJWT(token, blacklist)
		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "failed to check token blacklist")
	})

	t.Run("it fails when token is invalid", func(t *testing.T) {
		invalidToken := "invalid.jwt.token"

		// Create gomock blacklist
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		blacklist := NewMockTokenBlacklist(ctrl)
		blacklist.EXPECT().IsTokenBlacklisted(gomock.Any(), gomock.Any()).Return(false, nil).AnyTimes()

		// Validate invalid token
		claims, err := ValidateJWT(invalidToken, blacklist)
		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}

func TestValidateJWT_WithNilBlacklist(t *testing.T) {
	t.Parallel()

	originalSecret := os.Getenv("JWT_SECRET")
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
	}()

	t.Run("it works with nil blacklist", func(t *testing.T) {
		token, err := GenerateJWT(1, "Employee")
		require.NoError(t, err)

		claims, err := ValidateJWT(token, nil)
		require.NoError(t, err)
		assert.Equal(t, uint(1), claims.ID)
		assert.Equal(t, "Employee", claims.Role)
		assert.NotEmpty(t, claims.RegisteredClaims.ID)
	})
}
