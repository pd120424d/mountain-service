package model

import (
	"testing"

	"github.com/stretchr/testify/assert"

	urgencyV1 "github.com/pd120424d/mountain-service/api/contracts/urgency/v1"
)

func TestComputeSortPriority_FromDBStatusStrings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		statusStr          string
		assignedEmployeeID *uint
		expected           int
	}{
		{"open & unassigned maps to 1", string(urgencyV1.Open), nil, 1},
		{"open & assigned maps to 2", string(urgencyV1.Open), func() *uint { v := uint(42); return &v }(), 2},
		{"in_progress maps to 3", string(urgencyV1.InProgress), nil, 3},
		{"resolved maps to 4", string(urgencyV1.Resolved), nil, 4},
		{"closed maps to 5", string(urgencyV1.Closed), nil, 5},
		{"unknown maps to 6", "something_else", nil, 6},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			status := urgencyV1.UrgencyStatus(tc.statusStr)
			sp := ComputeSortPriority(status, tc.assignedEmployeeID)
			assert.Equal(t, tc.expected, sp)
		})
	}
}

func TestUpdateWithRequest_ThenComputeSortPriority(t *testing.T) {

	urg := &Urgency{
		Status: urgencyV1.UrgencyStatus(urgencyV1.Open),
	}

	req := &urgencyV1.UrgencyUpdateRequest{
		Status: urgencyV1.InProgress,
	}
	urg.UpdateWithRequest(req)

	sp := ComputeSortPriority(urgencyV1.UrgencyStatus(urg.Status), nil)
	assert.Equal(t, 3, sp)
}
