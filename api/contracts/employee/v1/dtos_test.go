package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmployeeCreateRequest_Validate(t *testing.T) {
	t.Parallel()

	t.Run("it returns no error for a valid request", func(t *testing.T) {
		req := &EmployeeCreateRequest{
			FirstName:      "John",
			LastName:       "Doe",
			Username:       "johndoe",
			Password:       "Pass123!",
			Email:          "john.doe@example.com",
			Gender:         "M",
			Phone:          "+1234567890",
			ProfilePicture: "https://example.com/profile.jpg",
			ProfileType:    "Medic",
		}

		err := req.Validate()
		assert.NoError(t, err)
	})

	t.Run("it returns an error for missing first name", func(t *testing.T) {
		req := &EmployeeCreateRequest{
			LastName:    "Doe",
			Username:    "johndoe",
			Password:    "Pass123!",
			Email:       "john.doe@example.com",
			Gender:      "M",
			Phone:       "+1234567890",
			ProfileType: "Medic",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "first name is required")
	})

	t.Run("it returns an error for missing last name", func(t *testing.T) {
		req := &EmployeeCreateRequest{
			FirstName:   "John",
			Username:    "johndoe",
			Password:    "Pass123!",
			Email:       "john.doe@example.com",
			Gender:      "M",
			Phone:       "+1234567890",
			ProfileType: "Medic",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "last name is required")
	})

	t.Run("it returns an error for missing username", func(t *testing.T) {
		req := &EmployeeCreateRequest{
			FirstName:   "John",
			LastName:    "Doe",
			Password:    "Pass123!",
			Email:       "john.doe@example.com",
			Gender:      "M",
			Phone:       "+1234567890",
			ProfileType: "Medic",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "username is required")
	})

	t.Run("it returns an error for missing password", func(t *testing.T) {
		req := &EmployeeCreateRequest{
			FirstName:   "John",
			LastName:    "Doe",
			Username:    "johndoe",
			Email:       "john.doe@example.com",
			Gender:      "M",
			Phone:       "+1234567890",
			ProfileType: "Medic",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password is required")
	})

	t.Run("it returns an error for invalid email", func(t *testing.T) {
		req := &EmployeeCreateRequest{
			FirstName:   "John",
			LastName:    "Doe",
			Username:    "johndoe",
			Password:    "Pass123!",
			Email:       "invalid-email",
			Gender:      "M",
			Phone:       "+1234567890",
			ProfileType: "Medic",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid email format")
	})

	t.Run("it returns an error for invalid gender", func(t *testing.T) {
		req := &EmployeeCreateRequest{
			FirstName:   "John",
			LastName:    "Doe",
			Username:    "johndoe",
			Password:    "Pass123!",
			Email:       "john.doe@example.com",
			Gender:      "X",
			Phone:       "+1234567890",
			ProfileType: "Medic",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "gender must be 'M' or 'F'")
	})

	t.Run("it returns an error for invalid phone", func(t *testing.T) {
		req := &EmployeeCreateRequest{
			FirstName:   "John",
			LastName:    "Doe",
			Username:    "johndoe",
			Password:    "Pass123!",
			Email:       "john.doe@example.com",
			Gender:      "M",
			Phone:       "123",
			ProfileType: "Medic",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must have 6-15 digits")
	})

	t.Run("it returns an error for invalid profile type", func(t *testing.T) {
		req := &EmployeeCreateRequest{
			FirstName:   "John",
			LastName:    "Doe",
			Username:    "johndoe",
			Password:    "Pass123!",
			Email:       "john.doe@example.com",
			Gender:      "M",
			Phone:       "+1234567890",
			ProfileType: "InvalidType",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "profile type must be one of: Medic, Technical, Administrator")
	})

	t.Run("it returns an error for missing profile type", func(t *testing.T) {
		req := &EmployeeCreateRequest{
			FirstName: "John",
			LastName:  "Doe",
			Username:  "johndoe",
			Password:  "Pass123!",
			Email:     "john.doe@example.com",
			Gender:    "M",
			Phone:     "+1234567890",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "profile type is required")
	})
}

func TestEmployeeUpdateRequest_Validate(t *testing.T) {
	t.Parallel()

	t.Run("it returns no error for a valid request", func(t *testing.T) {
		req := &EmployeeUpdateRequest{
			FirstName:      "John",
			LastName:       "Doe",
			Username:       "johndoe",
			Email:          "john.doe@example.com",
			Gender:         "M",
			Phone:          "+1234567890",
			ProfilePicture: "https://example.com/profile.jpg",
			ProfileType:    "Technical",
		}

		err := req.Validate()
		assert.NoError(t, err)
	})

	t.Run("it returns no error for empty fields", func(t *testing.T) {
		req := &EmployeeUpdateRequest{}

		err := req.Validate()
		assert.NoError(t, err)
	})

	t.Run("it returns an error for invalid email", func(t *testing.T) {
		req := &EmployeeUpdateRequest{
			Email: "invalid-email",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid email format")
	})

	t.Run("it returns an error for invalid gender", func(t *testing.T) {
		req := &EmployeeUpdateRequest{
			Gender: "X",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "gender must be 'M' or 'F'")
	})

	t.Run("it returns an error for invalid phone", func(t *testing.T) {
		req := &EmployeeUpdateRequest{
			Phone: "123",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must have 6-15 digits")
	})

	t.Run("it returns an error for invalid profile type", func(t *testing.T) {
		req := &EmployeeUpdateRequest{
			ProfileType: "InvalidType",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "profile type must be one of: Medic, Technical, Administrator")
	})
}

func TestEmployeeLogin_Validate(t *testing.T) {
	t.Parallel()

	t.Run("it returns no error for a valid request", func(t *testing.T) {
		req := &EmployeeLogin{
			Username: "johndoe",
			Password: "Pass123!",
		}

		err := req.Validate()
		assert.NoError(t, err)
	})

	t.Run("it returns an error for missing username", func(t *testing.T) {
		req := &EmployeeLogin{
			Password: "Pass123!",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "username is required")
	})

	t.Run("it returns an error for empty username", func(t *testing.T) {
		req := &EmployeeLogin{
			Username: "",
			Password: "Pass123!",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "username is required")
	})

	t.Run("it returns an error for whitespace-only username", func(t *testing.T) {
		req := &EmployeeLogin{
			Username: "   ",
			Password: "Pass123!",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "username is required")
	})

	t.Run("it returns an error for missing password", func(t *testing.T) {
		req := &EmployeeLogin{
			Username: "johndoe",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password is required")
	})

	t.Run("it returns an error for empty password", func(t *testing.T) {
		req := &EmployeeLogin{
			Username: "johndoe",
			Password: "",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password is required")
	})

	t.Run("it returns an error for whitespace-only password", func(t *testing.T) {
		req := &EmployeeLogin{
			Username: "johndoe",
			Password: "   ",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password is required")
	})

	// Password format validation tests removed - login validation only checks for required fields
	// Password format validation is handled during authentication, not at DTO level
	// This allows admin passwords to have different rules than employee passwords

	t.Run("it returns no error for simple password", func(t *testing.T) {
		req := &EmployeeLogin{
			Username: "johndoe",
			Password: "simple",
		}

		err := req.Validate()
		assert.NoError(t, err)
	})

	t.Run("it returns no error for admin-style password", func(t *testing.T) {
		req := &EmployeeLogin{
			Username: "admin",
			Password: "admin123",
		}

		err := req.Validate()
		assert.NoError(t, err)
	})

	t.Run("it returns multiple errors for invalid request", func(t *testing.T) {
		req := &EmployeeLogin{
			Username: "",
			Password: "",
		}

		err := req.Validate()
		assert.Error(t, err)
		// The error message should contain both validation errors
		errorMsg := err.Error()
		assert.Contains(t, errorMsg, "username is required")
		assert.Contains(t, errorMsg, "password is required")
	})
}
