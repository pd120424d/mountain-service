package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProfileTypeFromString(t *testing.T) {
	t.Run("it returns Medic when input is Medic", func(t *testing.T) {
		pt := ProfileTypeFromString("Medic")
		assert.Equal(t, Medic, pt)
	})

	t.Run("it returns Technical when input is Technical", func(t *testing.T) {
		pt := ProfileTypeFromString("Technical")
		assert.Equal(t, Technical, pt)
	})

	t.Run("it returns Administrator when input is Administrator", func(t *testing.T) {
		pt := ProfileTypeFromString("Administrator")
		assert.Equal(t, Administrator, pt)
	})

	t.Run("it returns empty string when input is invalid", func(t *testing.T) {
		pt := ProfileTypeFromString("Invalid")
		assert.Equal(t, ProfileType(""), pt)
	})
}

func TestProfileType_String(t *testing.T) {
	t.Run("it returns Medic when profile type is Medic", func(t *testing.T) {
		pt := Medic.String()
		assert.Equal(t, "Medic", pt)
	})

	t.Run("it returns Technical when profile type is Technical", func(t *testing.T) {
		pt := Technical.String()
		assert.Equal(t, "Technical", pt)
	})

	t.Run("it returns Administrator when profile type is Administrator", func(t *testing.T) {
		pt := Administrator.String()
		assert.Equal(t, "Administrator", pt)
	})

	t.Run("it returns Unknown when profile type is invalid", func(t *testing.T) {
		pt := ProfileType("Invalid").String()
		assert.Equal(t, "Unknown", pt)
	})
}

func TestProfileType_Valid(t *testing.T) {
	t.Run("it returns true when profile type is Medic", func(t *testing.T) {
		pt := Medic.Valid()
		assert.True(t, pt)
	})

	t.Run("it returns true when profile type is Technical", func(t *testing.T) {
		pt := Technical.Valid()
		assert.True(t, pt)
	})

	t.Run("it returns false when profile type is invalid", func(t *testing.T) {
		pt := ProfileType("Invalid").Valid()
		assert.False(t, pt)
	})
}

func TestEmployee_Role(t *testing.T) {
	employee := &Employee{
		ProfileType: Medic,
	}
	assert.Equal(t, "Medic", employee.Role())
}

func TestEmployee_UpdateResponseFromEmployee(t *testing.T) {
	employee := &Employee{
		ID:             1,
		Username:       "test-user",
		FirstName:      "Bruce",
		LastName:       "Lee",
		Gender:         "M",
		Phone:          "123456789",
		Email:          "test-user@example.com",
		ProfilePicture: "https://example.com/profile.jpg",
		ProfileType:    Medic,
	}

	response := employee.UpdateResponseFromEmployee()
	assert.Equal(t, EmployeeResponse{
		ID:             1,
		Username:       "test-user",
		FirstName:      "Bruce",
		LastName:       "Lee",
		Gender:         "M",
		Phone:          "123456789",
		Email:          "test-user@example.com",
		ProfilePicture: "https://example.com/profile.jpg",
		ProfileType:    "Medic",
	}, response)
}
