package utils

import (
	"fmt"
	"net/mail"
	"regexp"
	"strconv"
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

// ValidatePhone validates a phone number format
// Accepts various international formats including +, spaces, dashes, and parentheses
func ValidatePhone(phone string) error {
	if phone == "" {
		return fmt.Errorf("phone number is required")
	}

	// First check if there are any invalid characters (letters)
	if regexp.MustCompile(`[a-zA-Z]`).MatchString(phone) {
		return fmt.Errorf("invalid phone number format")
	}

	// Remove all non-digit characters except +
	cleaned := regexp.MustCompile(`[^\d+]`).ReplaceAllString(phone, "")

	// Check if it starts with + and has 7-15 digits after
	if strings.HasPrefix(cleaned, "+") {
		digits := cleaned[1:]
		if len(digits) < 7 || len(digits) > 15 {
			return fmt.Errorf("phone number must have 7-15 digits after country code")
		}
		// Validate all characters after + are digits
		if !regexp.MustCompile(`^\d+$`).MatchString(digits) {
			return fmt.Errorf("invalid phone number format")
		}
	} else {
		// Local format - should have 6-15 digits
		if len(cleaned) < 6 || len(cleaned) > 15 {
			return fmt.Errorf("phone number must have 6-15 digits")
		}
		// Validate all characters are digits
		if !regexp.MustCompile(`^\d+$`).MatchString(cleaned) {
			return fmt.Errorf("invalid phone number format")
		}
	}

	return nil
}

// ValidateOptionalPhone validates a phone number only if it's not empty
func ValidateOptionalPhone(phone string) error {
	if phone == "" {
		return nil
	}
	return ValidatePhone(phone)
}

// ValidateGender validates gender field (M or F)
func ValidateGender(gender string) error {
	if gender == "" {
		return fmt.Errorf("gender is required")
	}

	gender = strings.ToUpper(strings.TrimSpace(gender))
	if gender != "M" && gender != "F" {
		return fmt.Errorf("gender must be 'M' or 'F'")
	}

	return nil
}

// ValidateCoordinates validates GPS coordinates in format "N 43.401123 E 22.662756"
func ValidateCoordinates(coordinates string) error {
	if coordinates == "" {
		return fmt.Errorf("coordinates are required")
	}

	// Pattern for coordinates: N/S latitude E/W longitude
	// Allow for optional negative signs and decimal points
	pattern := `^[NS]\s*(-?\d+(?:\.\d+)?)\s*[EW]\s*(-?\d+(?:\.\d+)?)$`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(coordinates)
	if len(matches) != 3 {
		return fmt.Errorf("coordinates must be in format 'N 43.401123 E 22.662756'")
	}

	// Validate latitude range (-90 to 90)
	lat, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return fmt.Errorf("invalid latitude value")
	}
	if lat < -90 || lat > 90 {
		return fmt.Errorf("latitude must be between -90 and 90")
	}

	// Validate longitude range (-180 to 180)
	lng, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return fmt.Errorf("invalid longitude value")
	}
	if lng < -180 || lng > 180 {
		return fmt.Errorf("longitude must be between -180 and 180")
	}

	return nil
}

// ValidateOptionalCoordinates validates coordinates only if not empty
func ValidateOptionalCoordinates(coordinates string) error {
	if coordinates == "" {
		return nil
	}
	return ValidateCoordinates(coordinates)
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
