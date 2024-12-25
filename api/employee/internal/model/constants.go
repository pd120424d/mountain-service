package model

import "fmt"

const DateFormat = "2006-01-02"

var (
	ErrAlreadyAssigned = fmt.Errorf("already assigned")
	ErrCapacityReached = fmt.Errorf("capacity reached")
)
