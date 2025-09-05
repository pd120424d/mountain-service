package service

//go:generate mockgen -source=service.go -destination=service_gomock.go -package=service -imports=gomock=go.uber.org/mock/gomock

import (
	"context"
	"time"

	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
	"github.com/pd120424d/mountain-service/api/employee/internal/model"
)

// ShiftService handles all shift-related operations
type ShiftService interface {
	AssignShift(ctx context.Context, employeeID uint, req employeeV1.AssignShiftRequest) (*employeeV1.AssignShiftResponse, error)
	GetShifts(ctx context.Context, employeeID uint) ([]employeeV1.ShiftResponse, error)
	GetShiftsAvailability(ctx context.Context, employeeID uint, days int) (*employeeV1.ShiftAvailabilityResponse, error)
	RemoveShift(ctx context.Context, employeeID uint, req employeeV1.RemoveShiftRequest) error
	GetOnCallEmployees(ctx context.Context, currentTime time.Time, shiftBuffer time.Duration) ([]employeeV1.EmployeeResponse, error)
	GetShiftWarnings(ctx context.Context, employeeID uint) ([]string, error)

	GetAdminShiftsAvailability(ctx context.Context, days int) (*employeeV1.ShiftAvailabilityResponse, error)
}

// EmployeeService handles employee CRUD operations
type EmployeeService interface {
	RegisterEmployee(ctx context.Context, req employeeV1.EmployeeCreateRequest) (*employeeV1.EmployeeResponse, error)
	LoginEmployee(ctx context.Context, req employeeV1.EmployeeLogin) (string, error)
	LogoutEmployee(ctx context.Context, tokenID string, expiresAt time.Time) error
	ListEmployees(ctx context.Context) ([]employeeV1.EmployeeResponse, error)
	UpdateEmployee(ctx context.Context, employeeID uint, req employeeV1.EmployeeUpdateRequest) (*employeeV1.EmployeeResponse, error)
	DeleteEmployee(ctx context.Context, employeeID uint) error
	GetEmployeeByID(ctx context.Context, employeeID uint) (*model.Employee, error)
	GetEmployeeByUsername(ctx context.Context, username string) (*model.Employee, error)
	ResetAllData(ctx context.Context) error
}
