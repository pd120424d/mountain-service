package model

import (
	"testing"

	"github.com/stretchr/testify/assert"

	urgencyV1 "github.com/pd120424d/mountain-service/api/contracts/urgency/v1"
)

func TestUrgencyLevelFromString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected urgencyV1.UrgencyLevel
	}{
		{"Low", urgencyV1.Low},
		{"Medium", urgencyV1.Medium},
		{"High", urgencyV1.High},
		{"Critical", urgencyV1.Critical},
		{"invalid", urgencyV1.Medium},
		{"", urgencyV1.Medium},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := UrgencyLevelFromString(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}
