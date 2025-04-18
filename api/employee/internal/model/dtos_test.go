package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmployeeCreateRequest_ToString(t *testing.T) {
	req := &EmployeeCreateRequest{
		FirstName:      "Bruce",
		LastName:       "Lee",
		Username:       "test-user",
		Password:       "Pass123!",
		Email:          "test-user@example.com",
		Gender:         "M",
		Phone:          "123456789",
		ProfilePicture: "https://example.com/profile.jpg",
		ProfileType:    "Medic",
	}

	expected := "EmployeeCreateRequest { FirstName: Bruce, LastName: Lee, Username: test-user, Password: ********," +
		" Email: test-user@example.com, Gender: M, Phone: 123456789, ProfilePicture: https://example.com/profile.jpg, ProfileType: Medic }"
	assert.Equal(t, expected, req.ToString())
}

func TestEmployeeUpdateRequest_Validate(t *testing.T) {
	t.Run("it returns an error when email is invalid", func(t *testing.T) {
		req := &EmployeeUpdateRequest{
			Email: "invalid-email",
		}
		err := req.Validate()
		assert.Error(t, err)
		assert.EqualError(t, err, "mail: missing '@' or angle-addr")
	})

	t.Run("it returns no error when email is valid", func(t *testing.T) {
		req := &EmployeeUpdateRequest{
			Email: "test-user@example.com",
		}
		err := req.Validate()
		assert.NoError(t, err)
	})

	t.Run("it returns no error when email is empty", func(t *testing.T) {
		req := &EmployeeUpdateRequest{
			Email: "",
		}
		err := req.Validate()
		assert.NoError(t, err)
	})
}

func TestSanitizePassword(t *testing.T) {
	t.Run("it masks the password with asterisks", func(t *testing.T) {
		password := "Pass123!"
		sanitized := sanitizePassword(password)
		assert.Equal(t, "********", sanitized)
	})

	t.Run("it returns an empty string when password is empty", func(t *testing.T) {
		password := ""
		sanitized := sanitizePassword(password)
		assert.Equal(t, "", sanitized)
	})
}
