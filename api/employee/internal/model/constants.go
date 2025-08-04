package model

import "fmt"

var (
	ErrAlreadyAssigned = fmt.Errorf("already assigned")
	ErrCapacityReached = fmt.Errorf("capacity reached")
)

const (
	WarningInsufficientShifts   = "SHIFT_WARNINGS.INSUFFICIENT_SHIFTS"
	ErrorConsecutiveShiftsLimit = "SHIFT_ERRORS.CONSECUTIVE_SHIFTS_LIMIT"
)
