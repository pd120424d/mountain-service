package utils

import (
	"fmt"
	"net/mail"
	"strings"
)

// ValidateEmail validates an email address using Go's mail.ParseAddress
// Returns nil if email is valid, error otherwise
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// ValidateOptionalEmail validates an email address only if it's not empty
// Returns nil if email is empty or valid, error if invalid
func ValidateOptionalEmail(email string) error {
	if email == "" {
		return nil
	}
	_, err := mail.ParseAddress(email)
	return err
}

// ValidateRequiredField validates that a string field is not empty after trimming
// Returns error if field is empty, nil otherwise
func ValidateRequiredField(fieldValue, fieldName string) error {
	if strings.TrimSpace(fieldValue) == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	return nil
}

// SanitizePassword masks a password with asterisks for logging purposes
// Returns empty string if password is empty, otherwise returns "********"
func SanitizePassword(password string) string {
	if strings.TrimSpace(password) == "" {
		return ""
	}
	return "********"
}

// IsEmptyOrWhitespace checks if a string is empty or contains only whitespace
func IsEmptyOrWhitespace(s string) bool {
	return strings.TrimSpace(s) == ""
}
