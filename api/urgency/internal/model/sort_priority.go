package model

import (
	urgencyV1 "github.com/pd120424d/mountain-service/api/contracts/urgency/v1"
)

// ComputeSortPriority maps business rules to a small int for quick indexed sorting.
// 0: open & unassigned, 1: open & assigned, 2: in_progress, 3: resolved, 4: closed, 5: other
func ComputeSortPriority(status urgencyV1.UrgencyStatus, assignedEmployeeID *uint) int {
	switch status {
	case urgencyV1.UrgencyStatus("Open"), urgencyV1.UrgencyStatus("open"):
		if assignedEmployeeID == nil {
			return 0
		}
		return 1
	case urgencyV1.UrgencyStatus("InProgress"), urgencyV1.UrgencyStatus("in_progress"):
		return 2
	case urgencyV1.UrgencyStatus("Resolved"), urgencyV1.UrgencyStatus("resolved"):
		return 3
	case urgencyV1.UrgencyStatus("Closed"), urgencyV1.UrgencyStatus("closed"):
		return 4
	default:
		return 5
	}
}

