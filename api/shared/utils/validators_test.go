package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	t.Parallel()

	t.Run("it returns an error when email is empty", func(t *testing.T) {
		err := ValidateEmail("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email is required")
	})

	t.Run("it returns an error when email is invalid", func(t *testing.T) {
		err := ValidateEmail("invalid-email")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid email format")
	})

	t.Run("it returns no error when email is valid", func(t *testing.T) {
		err := ValidateEmail("test@example.com")
		assert.NoError(t, err)
	})

	t.Run("it returns no error for complex valid email", func(t *testing.T) {
		err := ValidateEmail("user.name+tag@example.co.uk")
		assert.NoError(t, err)
	})
}

func TestValidateOptionalEmail(t *testing.T) {
	t.Run("it returns no error when email is empty", func(t *testing.T) {
		err := ValidateOptionalEmail("")
		assert.NoError(t, err)
	})

	t.Run("it returns an error when email is invalid", func(t *testing.T) {
		err := ValidateOptionalEmail("invalid-email")
		assert.Error(t, err)
		assert.EqualError(t, err, "mail: missing '@' or angle-addr")
	})

	t.Run("it returns no error when email is valid", func(t *testing.T) {
		err := ValidateOptionalEmail("test@example.com")
		assert.NoError(t, err)
	})
}

func TestValidateRequiredField(t *testing.T) {
	t.Run("it returns an error when field is empty", func(t *testing.T) {
		err := ValidateRequiredField("", "name")
		assert.Error(t, err)
		assert.EqualError(t, err, "name is required")
	})

	t.Run("it returns an error when field contains only whitespace", func(t *testing.T) {
		err := ValidateRequiredField("   ", "description")
		assert.Error(t, err)
		assert.EqualError(t, err, "description is required")
	})

	t.Run("it returns no error when field has valid content", func(t *testing.T) {
		err := ValidateRequiredField("valid content", "name")
		assert.NoError(t, err)
	})

	t.Run("it returns no error when field has content with surrounding whitespace", func(t *testing.T) {
		err := ValidateRequiredField("  valid content  ", "name")
		assert.NoError(t, err)
	})
}

func TestSanitizePassword(t *testing.T) {
	t.Run("it masks the password with asterisks", func(t *testing.T) {
		password := "Pass123!"
		sanitized := SanitizePassword(password)
		assert.Equal(t, "********", sanitized)
	})

	t.Run("it returns an empty string when password is empty", func(t *testing.T) {
		password := ""
		sanitized := SanitizePassword(password)
		assert.Equal(t, "", sanitized)
	})

	t.Run("it returns an empty string when password contains only whitespace", func(t *testing.T) {
		password := "   "
		sanitized := SanitizePassword(password)
		assert.Equal(t, "", sanitized)
	})

	t.Run("it masks short passwords", func(t *testing.T) {
		password := "a"
		sanitized := SanitizePassword(password)
		assert.Equal(t, "********", sanitized)
	})

	t.Run("it masks long passwords", func(t *testing.T) {
		password := "ThisIsAVeryLongPasswordThatShouldStillBeMasked"
		sanitized := SanitizePassword(password)
		assert.Equal(t, "********", sanitized)
	})
}

func TestIsEmptyOrWhitespace(t *testing.T) {
	t.Run("it returns true for empty string", func(t *testing.T) {
		result := IsEmptyOrWhitespace("")
		assert.True(t, result)
	})

	t.Run("it returns true for whitespace only", func(t *testing.T) {
		result := IsEmptyOrWhitespace("   ")
		assert.True(t, result)
	})

	t.Run("it returns true for tabs and newlines", func(t *testing.T) {
		result := IsEmptyOrWhitespace("\t\n\r ")
		assert.True(t, result)
	})

	t.Run("it returns false for string with content", func(t *testing.T) {
		result := IsEmptyOrWhitespace("content")
		assert.False(t, result)
	})

	t.Run("it returns false for string with content and whitespace", func(t *testing.T) {
		result := IsEmptyOrWhitespace("  content  ")
		assert.False(t, result)
	})
}
