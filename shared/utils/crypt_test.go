package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePassword(t *testing.T) {
	t.Run("it returns an error when password is too short", func(t *testing.T) {
		err := ValidatePassword("Abc1!")
		assert.EqualError(t, err, ErrPasswordLength)
	})

	t.Run("it returns an error when password does not contain uppercase letter", func(t *testing.T) {
		err := ValidatePassword("abc123!")
		assert.EqualError(t, err, ErrMissingUppercase.Error())
	})

	t.Run("it returns an error when password does not contain enough lowercase letters", func(t *testing.T) {
		err := ValidatePassword("ABC1de!")
		assert.EqualError(t, err, ErrInsufficientLowercase.Error())
	})

	t.Run("it returns an error when password does not contain a digit", func(t *testing.T) {
		err := ValidatePassword("Abcdef!")
		assert.EqualError(t, err, ErrMissingDigit.Error())
	})

	t.Run("it returns an error when password does not contain a special character", func(t *testing.T) {
		err := ValidatePassword("Abcdef1")
		assert.EqualError(t, err, ErrMissingSpecialCharacter.Error())
	})

	t.Run("it returns an error when password does not start with a letter", func(t *testing.T) {
		err := ValidatePassword("1Abcdef!")
		assert.EqualError(t, err, ErrPasswordMustStartWithLetter.Error())
	})

	t.Run("it does not return an error when password is valid", func(t *testing.T) {
		err := ValidatePassword("Abcd123!")
		assert.NoError(t, err)
	})
}
