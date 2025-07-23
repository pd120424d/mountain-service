package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewServiceAuth(t *testing.T) {
	t.Parallel()

	t.Run("it generates a valid token", func(t *testing.T) {
		serviceAuth := NewServiceAuth(ServiceAuthConfig{
			Secret:      "test-secret",
			ServiceName: "test-service",
			TokenTTL:    1 * time.Hour,
		})
		token, err := serviceAuth.GenerateToken()
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("it validates a valid token", func(t *testing.T) {
		serviceAuth := NewServiceAuth(ServiceAuthConfig{
			Secret:      "test-secret",
			ServiceName: "test-service",
			TokenTTL:    1 * time.Hour,
		})
		token, err := serviceAuth.GenerateToken()
		assert.NoError(t, err)

		claims, err := serviceAuth.ValidateToken(token)
		assert.NoError(t, err)
		assert.Equal(t, "test-service", claims.ServiceName)
	})

	t.Run("it returns the Authorization header value", func(t *testing.T) {
		serviceAuth := NewServiceAuth(ServiceAuthConfig{
			Secret:      "test-secret",
			ServiceName: "test-service",
			TokenTTL:    1 * time.Hour,
		})
		token, err := serviceAuth.GenerateToken()
		assert.NoError(t, err)

		authHeader, err := serviceAuth.GetAuthHeader()
		assert.NoError(t, err)
		assert.Equal(t, "Bearer "+token, authHeader)
	})
}
