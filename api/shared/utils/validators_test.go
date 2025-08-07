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

func TestValidatePhone(t *testing.T) {
	t.Parallel()

	t.Run("it returns an error when phone is empty", func(t *testing.T) {
		err := ValidatePhone("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "phone number is required")
	})

	t.Run("it returns no error for valid international phone", func(t *testing.T) {
		err := ValidatePhone("+1234567890")
		assert.NoError(t, err)
	})

	t.Run("it returns no error for valid local phone", func(t *testing.T) {
		err := ValidatePhone("1234567890")
		assert.NoError(t, err)
	})

	t.Run("it returns no error for phone with formatting", func(t *testing.T) {
		err := ValidatePhone("+1 (234) 567-8900")
		assert.NoError(t, err)
	})

	t.Run("it returns an error for too short phone", func(t *testing.T) {
		err := ValidatePhone("12345")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must have 6-15 digits")
	})

	t.Run("it returns an error for too long phone", func(t *testing.T) {
		err := ValidatePhone("1234567890123456")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must have 6-15 digits")
	})

	t.Run("it returns an error for invalid characters", func(t *testing.T) {
		err := ValidatePhone("123abc456")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid phone number format")
	})
}

func TestValidateOptionalPhone(t *testing.T) {
	t.Run("it returns no error when phone is empty", func(t *testing.T) {
		err := ValidateOptionalPhone("")
		assert.NoError(t, err)
	})

	t.Run("it returns an error when phone is invalid", func(t *testing.T) {
		err := ValidateOptionalPhone("123")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must have 6-15 digits")
	})

	t.Run("it returns no error when phone is valid", func(t *testing.T) {
		err := ValidateOptionalPhone("+1234567890")
		assert.NoError(t, err)
	})
}

func TestValidateGender(t *testing.T) {
	t.Parallel()

	t.Run("it returns an error when gender is empty", func(t *testing.T) {
		err := ValidateGender("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "gender is required")
	})

	t.Run("it returns no error for valid gender M", func(t *testing.T) {
		err := ValidateGender("M")
		assert.NoError(t, err)
	})

	t.Run("it returns no error for valid gender F", func(t *testing.T) {
		err := ValidateGender("F")
		assert.NoError(t, err)
	})

	t.Run("it returns no error for lowercase gender", func(t *testing.T) {
		err := ValidateGender("m")
		assert.NoError(t, err)
	})

	t.Run("it returns an error for invalid gender", func(t *testing.T) {
		err := ValidateGender("X")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "gender must be 'M' or 'F'")
	})
}

func TestValidateCoordinates(t *testing.T) {
	t.Parallel()

	t.Run("it returns an error when coordinates are empty", func(t *testing.T) {
		err := ValidateCoordinates("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "coordinates are required")
	})

	t.Run("it returns no error for valid coordinates", func(t *testing.T) {
		err := ValidateCoordinates("N 43.401123 E 22.662756")
		assert.NoError(t, err)
	})

	t.Run("it returns no error for coordinates with S and W", func(t *testing.T) {
		err := ValidateCoordinates("S 43.401123 W 22.662756")
		assert.NoError(t, err)
	})

	t.Run("it returns no error for coordinates with negative values", func(t *testing.T) {
		err := ValidateCoordinates("N -43.401123 E -22.662756")
		assert.NoError(t, err)
	})

	t.Run("it returns an error for invalid format", func(t *testing.T) {
		err := ValidateCoordinates("43.401123, 22.662756")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "coordinates must be in format")
	})

	t.Run("it returns an error for latitude out of range", func(t *testing.T) {
		err := ValidateCoordinates("N 91.0 E 22.662756")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "latitude must be between -90 and 90")
	})

	t.Run("it returns an error for longitude out of range", func(t *testing.T) {
		err := ValidateCoordinates("N 43.401123 E 181.0")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "longitude must be between -180 and 180")
	})

	t.Run("it returns an error for invalid latitude value", func(t *testing.T) {
		err := ValidateCoordinates("N abc E 22.662756")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "coordinates must be in format")
	})
}

func TestValidateOptionalCoordinates(t *testing.T) {
	t.Run("it returns no error when coordinates are empty", func(t *testing.T) {
		err := ValidateOptionalCoordinates("")
		assert.NoError(t, err)
	})

	t.Run("it returns an error when coordinates are invalid", func(t *testing.T) {
		err := ValidateOptionalCoordinates("invalid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "coordinates must be in format")
	})

	t.Run("it returns no error when coordinates are valid", func(t *testing.T) {
		err := ValidateOptionalCoordinates("N 43.401123 E 22.662756")
		assert.NoError(t, err)
	})
}
