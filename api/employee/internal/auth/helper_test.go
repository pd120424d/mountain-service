package auth

import (
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
