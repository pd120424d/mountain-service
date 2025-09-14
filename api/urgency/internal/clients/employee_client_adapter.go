package clients

import (
	"context"
	"time"

	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
	s2semployee "github.com/pd120424d/mountain-service/api/shared/s2s/employee"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

// NewEmployeeClientFromS2S adapts a shared s2s employee client to this package's EmployeeClient interface.
func NewEmployeeClientFromS2S(inner s2semployee.Client, logger utils.Logger) EmployeeClient {
	return &s2sEmployeeAdapter{inner: inner, logger: logger.WithName("employeeS2SAdapter")}
}

type s2sEmployeeAdapter struct {
	inner  s2semployee.Client
	logger utils.Logger
}

func (a *s2sEmployeeAdapter) GetEmployeeByID(ctx context.Context, employeeID uint) (*employeeV1.EmployeeResponse, error) {
	return a.inner.GetEmployeeByID(ctx, employeeID)
}

func (a *s2sEmployeeAdapter) GetAllEmployees(ctx context.Context) ([]employeeV1.EmployeeResponse, error) {
	return a.inner.GetAllEmployees(ctx)
}

func (a *s2sEmployeeAdapter) GetOnCallEmployees(ctx context.Context, shiftBuffer time.Duration) ([]employeeV1.EmployeeResponse, error) {
	return a.inner.GetOnCallEmployees(ctx, shiftBuffer)
}

func (a *s2sEmployeeAdapter) CheckActiveEmergencies(ctx context.Context, employeeID uint) (bool, error) {
	return a.inner.CheckActiveEmergencies(ctx, employeeID)
}

