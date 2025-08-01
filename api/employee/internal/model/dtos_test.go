package model

import (
	"testing"

	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/stretchr/testify/assert"
)

func TestEmployeeCreateRequest_ToString(t *testing.T) {
	req := &employeeV1.EmployeeCreateRequest{
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
		req := &employeeV1.EmployeeUpdateRequest{
			Email: "invalid-email",
		}
		err := utils.ValidateOptionalEmail(req.Email)
		assert.Error(t, err)
		assert.EqualError(t, err, "mail: missing '@' or angle-addr")
	})

	t.Run("it returns no error when email is valid", func(t *testing.T) {
		req := &employeeV1.EmployeeUpdateRequest{
			Email: "test-user@example.com",
		}
		err := utils.ValidateOptionalEmail(req.Email)
		assert.NoError(t, err)
	})

	t.Run("it returns no error when email is empty", func(t *testing.T) {
		req := &employeeV1.EmployeeUpdateRequest{
			Email: "",
		}
		err := utils.ValidateOptionalEmail(req.Email)
		assert.NoError(t, err)
	})
}

// TestSanitizePassword moved to shared/utils/validators_test.go
