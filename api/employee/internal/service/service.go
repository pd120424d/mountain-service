package service

//go:generate mockgen -source=service.go -destination=service_gomock.go -package=service -imports=gomock=go.uber.org/mock/gomock

import (
	"time"

	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
	"github.com/pd120424d/mountain-service/api/employee/internal/model"
)

// ShiftService handles all shift-related operations
type ShiftService interface {
	AssignShift(employeeID uint, req employeeV1.AssignShiftRequest) (*employeeV1.AssignShiftResponse, error)
	GetShifts(employeeID uint) ([]employeeV1.ShiftResponse, error)
	GetShiftsAvailability(employeeID uint, days int) (*employeeV1.ShiftAvailabilityResponse, error)
	RemoveShift(employeeID uint, req employeeV1.RemoveShiftRequest) error
	GetOnCallEmployees(currentTime time.Time, shiftBuffer time.Duration) ([]employeeV1.EmployeeResponse, error)
	GetShiftWarnings(employeeID uint) ([]string, error)
}

// EmployeeService handles employee CRUD operations
type EmployeeService interface {
	RegisterEmployee(req employeeV1.EmployeeCreateRequest) (*employeeV1.EmployeeResponse, error)
	LoginEmployee(req employeeV1.EmployeeLogin) (string, error)
	ListEmployees() ([]employeeV1.EmployeeResponse, error)
	UpdateEmployee(employeeID uint, req employeeV1.EmployeeUpdateRequest) (*employeeV1.EmployeeResponse, error)
	DeleteEmployee(employeeID uint) error
	GetEmployeeByID(employeeID uint) (*model.Employee, error)
	GetEmployeeByUsername(username string) (*model.Employee, error)
	ResetAllData() error
}
