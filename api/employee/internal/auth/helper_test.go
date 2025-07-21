package auth

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	t.Run("It generates a valid JWT with employee id and role", func(t *testing.T) {
		token, err := GenerateJWT(1, "Medic")

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}

func TestValidateJWT(t *testing.T) {
	t.Run("It returns claims when token is valid", func(t *testing.T) {
		token, err := GenerateJWT(1, "Medic")
		assert.NoError(t, err)

		claims, err := ValidateJWT(token)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), claims.ID)
		assert.Equal(t, "Medic", claims.Role)
	})
}

func TestIsAdminLogin(t *testing.T) {
	t.Run("It returns true for admin username", func(t *testing.T) {
		result := IsAdminLogin("admin")
		assert.True(t, result)
	})

	t.Run("It returns false for non-admin username", func(t *testing.T) {
		result := IsAdminLogin("user")
		assert.False(t, result)
	})
}

func TestValidateAdminPassword(t *testing.T) {
	t.Run("It validates password from environment variable", func(t *testing.T) {
		os.Setenv("ADMIN_PASSWORD", "testpass")
		defer os.Unsetenv("ADMIN_PASSWORD")

		result := ValidateAdminPassword("testpass")
		assert.True(t, result)

		result = ValidateAdminPassword("wrongpass")
		assert.False(t, result)
	})

	t.Run("It returns false when no password is set", func(t *testing.T) {
		os.Unsetenv("ADMIN_PASSWORD")
		os.Unsetenv("ADMIN_PASSWORD_FILE")

		result := ValidateAdminPassword("anypass")
		assert.False(t, result)
	})
}

func TestGenerateAdminJWT(t *testing.T) {
	t.Run("It generates a valid admin JWT", func(t *testing.T) {
		token, err := GenerateAdminJWT()
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		claims, err := ValidateJWT(token)
		assert.NoError(t, err)
		assert.Equal(t, uint(0), claims.ID)
		assert.Equal(t, "Administrator", claims.Role)
	})
}

func TestHashPassword(t *testing.T) {
	t.Run("It hashes the password", func(t *testing.T) {
		hashedPassword, err := HashPassword("password")
		assert.NoError(t, err)
		assert.NotEqual(t, "password", hashedPassword)
	})
}

func TestCheckPassword(t *testing.T) {
	t.Run("It returns true when password matches", func(t *testing.T) {
		hashedPassword, err := HashPassword("password")
		assert.NoError(t, err)

		match := CheckPassword(hashedPassword, "password")
		assert.True(t, match)
	})
}
