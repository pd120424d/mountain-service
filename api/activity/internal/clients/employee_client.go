package clients

//go:generate mockgen -source=employee_client.go -destination=employee_client_gomock.go -package=clients mountain_service/activity/internal/clients -imports=gomock=go.uber.org/mock/gomock -typed

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
	"github.com/pd120424d/mountain-service/api/shared/auth"
	"github.com/pd120424d/mountain-service/api/shared/client"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

type EmployeeClient interface {
	GetEmployeeByID(ctx context.Context, employeeID uint) (*employeeV1.EmployeeResponse, error)
}

type employeeHTTPClient interface {
	Get(ctx context.Context, endpoint string) (*http.Response, error)
}

type employeeClient struct {
	httpClient employeeHTTPClient
	logger     utils.Logger
}

type EmployeeClientConfig struct {
	BaseURL     string
	ServiceAuth auth.ServiceAuth
	Logger      utils.Logger
	Timeout     time.Duration
}

func NewEmployeeClient(config EmployeeClientConfig) EmployeeClient {
	httpClient := client.NewHTTPClient(client.HTTPClientConfig{
		BaseURL:     config.BaseURL,
		ServiceAuth: config.ServiceAuth,
		Logger:      config.Logger,
		Timeout:     config.Timeout,
	})
	return &employeeClient{httpClient: httpClient, logger: config.Logger.WithName("employeeClient")}
}

func (c *employeeClient) GetEmployeeByID(ctx context.Context, employeeID uint) (*employeeV1.EmployeeResponse, error) {
	log := c.logger.WithContext(ctx)
	endpoint := fmt.Sprintf("/api/v1/service/employees/%d", employeeID)
	log.Debugf("employee_client.get_by_id start id=%d endpoint=%s", employeeID, endpoint)

	resp, err := c.httpClient.Get(ctx, endpoint)
	if err != nil {
		log.Errorf("employee_client.get_by_id http_error id=%d err=%v", employeeID, err)
		return nil, fmt.Errorf("failed to call employee service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warnf("employee_client.get_by_id non_200 id=%d status=%d", employeeID, resp.StatusCode)
		return nil, fmt.Errorf("employee service returned status %d", resp.StatusCode)
	}

	var emp employeeV1.EmployeeResponse
	if err := json.NewDecoder(resp.Body).Decode(&emp); err != nil {
		log.Errorf("employee_client.get_by_id decode_error id=%d err=%v", employeeID, err)
		return nil, fmt.Errorf("failed to decode employee response: %w", err)
	}

	log.Debugf("employee_client.get_by_id success id=%d name=%s %s", employeeID, emp.FirstName, emp.LastName)
	return &emp, nil
}
