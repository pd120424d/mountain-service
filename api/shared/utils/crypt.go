package utils

import (
	"errors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const (
	ErrPasswordLength          = "password must be between 6 and 10 characters long"
	ErrPasswordUppercase       = "password must contain at least one uppercase letter"
	ErrPasswordLowercase       = "password must contain at least three lowercase letters"
	ErrPasswordDigit           = "password must contain at least one digit"
	ErrPasswordSpecial         = "password must contain at least one special character"
	ErrPasswordStartWithLetter = "password must start with a letter"
)

var (
	ErrInvalidPasswordLength       = errors.New(ErrPasswordLength)
	ErrMissingUppercase            = errors.New(ErrPasswordUppercase)
	ErrInsufficientLowercase       = errors.New(ErrPasswordLowercase)
	ErrMissingDigit                = errors.New(ErrPasswordDigit)
	ErrMissingSpecialCharacter     = errors.New(ErrPasswordSpecial)
	ErrPasswordMustStartWithLetter = errors.New(ErrPasswordStartWithLetter)
)

// HashPassword hashes the given password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// ValidatePassword checks if the given password meets the following criteria:
// 1. The password must be between 6 and 20 characters long.
// 2. The password must contain at least one uppercase letter.
// 3. The password must contain at least three lowercase letters.
// 4. The password must contain at least one digit.
// 5. The password must contain at least one special character.
func ValidatePassword(password string) error {
	if len(password) < 6 || len(password) > 20 {
		return ErrInvalidPasswordLength
	}

	var hasUpper, hasNumber, hasSpecial bool
	var lowerCount int

	for i, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			lowerCount++
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
		if i == 0 && !unicode.IsLetter(char) {
			return ErrPasswordMustStartWithLetter
		}
	}

	if !hasUpper {
		return ErrMissingUppercase
	}
	if lowerCount < 3 {
		return ErrInsufficientLowercase
	}
	if !hasNumber {
		return ErrMissingDigit
	}
	if !hasSpecial {
		return ErrMissingSpecialCharacter
	}

	return nil
}
