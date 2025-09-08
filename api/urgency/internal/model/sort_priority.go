package model

import (
	urgencyV1 "github.com/pd120424d/mountain-service/api/contracts/urgency/v1"
)

// ComputeSortPriority maps business rules to a small int for quick indexed sorting.
// NOTE: shifted by +1 to avoid zero, so GORM doesn’t omit the column when a DB default exists.
// 1: open & unassigned, 2: open & assigned, 3: in_progress, 4: resolved, 5: closed, 6: other
func ComputeSortPriority(status urgencyV1.UrgencyStatus, assignedEmployeeID *uint) int {
	switch status {
	case urgencyV1.Open:
		if assignedEmployeeID == nil {
			return 1
		}
		return 2
	case urgencyV1.InProgress:
		return 3
	case urgencyV1.Resolved:
		return 4
	case urgencyV1.Closed:
		return 5
	default:
		return 6
	}
}
