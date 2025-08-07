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

	t.Run("it returns an error for password too short", func(t *testing.T) {
		req := &EmployeeLogin{
			Username: "johndoe",
			Password: "Abc1!",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password must be between 6 and 10 characters long")
	})

	t.Run("it returns an error for password too long", func(t *testing.T) {
		req := &EmployeeLogin{
			Username: "johndoe",
			Password: "Abcdefgh123!",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password must be between 6 and 10 characters long")
	})

	t.Run("it returns an error for password without uppercase", func(t *testing.T) {
		req := &EmployeeLogin{
			Username: "johndoe",
			Password: "abcd123!",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password must contain at least one uppercase letter")
	})

	t.Run("it returns an error for password without enough lowercase", func(t *testing.T) {
		req := &EmployeeLogin{
			Username: "johndoe",
			Password: "AB12de!",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password must contain at least three lowercase letters")
	})

	t.Run("it returns an error for password without digit", func(t *testing.T) {
		req := &EmployeeLogin{
			Username: "johndoe",
			Password: "Abcdef!",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password must contain at least one digit")
	})

	t.Run("it returns an error for password without special character", func(t *testing.T) {
		req := &EmployeeLogin{
			Username: "johndoe",
			Password: "Abcdef1",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password must contain at least one special character")
	})

	t.Run("it returns an error for password not starting with letter", func(t *testing.T) {
		req := &EmployeeLogin{
			Username: "johndoe",
			Password: "1Abcde!",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password must start with a letter")
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
