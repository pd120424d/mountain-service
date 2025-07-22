package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUrgencyLevel_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		level    UrgencyLevel
		expected string
	}{
		{Low, "Low"},
		{Medium, "Medium"},
		{High, "High"},
		{Critical, "Critical"},
		{"invalid", "invalid"},
	}

	for _, test := range tests {
		t.Run(string(test.level), func(t *testing.T) {
			result := test.level.String()
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestUrgencyLevel_Valid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		level    UrgencyLevel
		expected bool
	}{
		{Low, true},
		{Medium, true},
		{High, true},
		{Critical, true},
		{UrgencyLevel("invalid"), false},
		{UrgencyLevel(""), false},
	}

	for _, test := range tests {
		t.Run(string(test.level), func(t *testing.T) {
			result := test.level.Valid()
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestUrgencyLevelFromString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected UrgencyLevel
	}{
		{"Low", Low},
		{"Medium", Medium},
		{"High", High},
		{"Critical", Critical},
		{"invalid", Medium},
		{"", Medium},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := UrgencyLevelFromString(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestUrgency_ToResponse(t *testing.T) {
	urgency := &Urgency{
		ID:           1,
		Name:         "Test Urgency",
		Email:        "test@example.com",
		ContactPhone: "123456789",
		Location:     "N 43.401123 E 22.662756",
		Description:  "Test description",
		Level:        High,
		Status:       "Open",
	}
	urgency.CreatedAt = time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	urgency.UpdatedAt = time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC)

	response := urgency.ToResponse()

	expected := UrgencyResponse{
		ID:           1,
		Name:         "Test Urgency",
		Email:        "test@example.com",
		ContactPhone: "123456789",
		Location:     "N 43.401123 E 22.662756",
		Description:  "Test description",
		Level:        "High",
		Status:       "Open",
		CreatedAt:    "2023-01-01T12:00:00Z",
		UpdatedAt:    "2023-01-02T12:00:00Z",
	}

	assert.Equal(t, expected, response)
}
