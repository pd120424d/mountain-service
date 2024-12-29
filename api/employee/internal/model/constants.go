package model

import "fmt"

var (
	ErrAlreadyAssigned = fmt.Errorf("already assigned")
	ErrCapacityReached = fmt.Errorf("capacity reached")
)
