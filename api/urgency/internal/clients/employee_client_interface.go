package clients

import (
	"context"
	"time"

	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
)

// EmployeeClient describes operations required from the employee service.
type EmployeeClient interface {
	GetOnCallEmployees(ctx context.Context, shiftBuffer time.Duration) ([]employeeV1.EmployeeResponse, error)
	GetAllEmployees(ctx context.Context) ([]employeeV1.EmployeeResponse, error)
	GetEmployeeByID(ctx context.Context, employeeID uint) (*employeeV1.EmployeeResponse, error)
	CheckActiveEmergencies(ctx context.Context, employeeID uint) (bool, error)
}

