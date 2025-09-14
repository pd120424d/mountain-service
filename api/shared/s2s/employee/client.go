package employee

//go:generate mockgen -source=client.go -destination=client_gomock.go -package=employee shared/s2s/employee -imports=gomock=go.uber.org/mock/gomock -typed

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
	"github.com/pd120424d/mountain-service/api/shared/auth"
	"github.com/pd120424d/mountain-service/api/shared/client"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

// Client defines the S2S client for the employee service.
type Client interface {
	GetEmployeeByID(ctx context.Context, employeeID uint) (*employeeV1.EmployeeResponse, error)
	GetAllEmployees(ctx context.Context) ([]employeeV1.EmployeeResponse, error)
	GetOnCallEmployees(ctx context.Context, shiftBuffer time.Duration) ([]employeeV1.EmployeeResponse, error)
	CheckActiveEmergencies(ctx context.Context, employeeID uint) (bool, error)
}

type httpClient interface {
	Get(ctx context.Context, endpoint string) (*http.Response, error)
}

type Config struct {
	BaseURL     string
	ServiceAuth auth.ServiceAuth
	Logger      utils.Logger
	Timeout     time.Duration
	MaxRetries  int
}

type clientImpl struct {
	http       httpClient
	logger     utils.Logger
	maxRetries int
}

// New constructs a client using the shared HTTP client.
func New(cfg Config) Client {
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	if cfg.MaxRetries <= 0 {
		cfg.MaxRetries = 2
	}
	h := client.NewHTTPClient(client.HTTPClientConfig{
		BaseURL:     cfg.BaseURL,
		Timeout:     cfg.Timeout,
		ServiceAuth: cfg.ServiceAuth,
		Logger:      cfg.Logger,
	})
	return &clientImpl{http: h, logger: cfg.Logger.WithName("employeeS2S"), maxRetries: cfg.MaxRetries}
}

// NewFromEnv constructs a Client using EMPLOYEE_SERVICE_URL or a sane in-cluster default.
func NewFromEnv(logger utils.Logger, sa auth.ServiceAuth) Client {
	base := os.Getenv("EMPLOYEE_SERVICE_URL")
	if base == "" {
		base = "http://employee-service:8082"
	}
	return New(Config{BaseURL: base, ServiceAuth: sa, Logger: logger, Timeout: 30 * time.Second})
}

func (c *clientImpl) retryGet(ctx context.Context, endpoint string) (*http.Response, error) {
	var lastErr error
	backoff := 100 * time.Millisecond
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		resp, err := c.http.Get(ctx, endpoint)
		if err == nil {
			// Retry on 429/5xx
			if resp.StatusCode == http.StatusTooManyRequests || (resp.StatusCode >= 500 && resp.StatusCode <= 599) {
				lastErr = fmt.Errorf("server status %d", resp.StatusCode)
				resp.Body.Close()
			} else {
				return resp, nil
			}
		} else {
			// Retry on transient network errors
			if ne, ok := err.(net.Error); ok && (ne.Timeout() || ne.Temporary()) {
				lastErr = err
			} else if err == context.DeadlineExceeded || err == context.Canceled {
				return nil, err
			} else {
				lastErr = err
			}
		}
		if attempt == c.maxRetries {
			break
		}
		select {
		case <-time.After(backoff):
			backoff = time.Duration(float64(backoff) * 1.8)
			continue
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return nil, lastErr
}

func (c *clientImpl) GetEmployeeByID(ctx context.Context, employeeID uint) (*employeeV1.EmployeeResponse, error) {
	log := c.logger.WithContext(ctx)
	endpoint := fmt.Sprintf("/api/v1/service/employees/%d", employeeID)
	resp, err := c.retryGet(ctx, endpoint)
	if err != nil {
		log.Errorf("employee.get_by_id http_error id=%d err=%v", employeeID, err)
		return nil, fmt.Errorf("failed to call employee service: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		log.Warnf("employee.get_by_id not_found id=%d", employeeID)
		return nil, fmt.Errorf("employee %d not found", employeeID)
	}
	if resp.StatusCode != http.StatusOK {
		log.Errorf("employee.get_by_id non_200 id=%d status=%d", employeeID, resp.StatusCode)
		return nil, fmt.Errorf("employee service returned status %d", resp.StatusCode)
	}
	var out employeeV1.EmployeeResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		log.Errorf("employee.get_by_id decode_error id=%d err=%v", employeeID, err)
		return nil, fmt.Errorf("failed to decode employee response: %w", err)
	}
	return &out, nil
}

func (c *clientImpl) GetAllEmployees(ctx context.Context) ([]employeeV1.EmployeeResponse, error) {
	log := c.logger.WithContext(ctx)
	resp, err := c.retryGet(ctx, "/api/v1/employees")
	if err != nil {
		log.Errorf("employee.get_all http_error err=%v", err)
		return nil, fmt.Errorf("failed to get all employees: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Errorf("employee.get_all non_200 status=%d", resp.StatusCode)
		return nil, fmt.Errorf("employee service returned status %d", resp.StatusCode)
	}
	var res employeeV1.AllEmployeesResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Errorf("employee.get_all decode_error err=%v", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return res.Employees, nil
}

func (c *clientImpl) GetOnCallEmployees(ctx context.Context, shiftBuffer time.Duration) ([]employeeV1.EmployeeResponse, error) {
	log := c.logger.WithContext(ctx)
	endpoint := "/api/v1/employees/on-call"
	if shiftBuffer > 0 {
		endpoint += fmt.Sprintf("?shift_buffer=%s", shiftBuffer.String())
	}
	resp, err := c.retryGet(ctx, endpoint)
	if err != nil {
		log.Errorf("employee.on_call http_error err=%v", err)
		return nil, fmt.Errorf("failed to get on-call employees: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Errorf("employee.on_call non_200 status=%d", resp.StatusCode)
		return nil, fmt.Errorf("employee service returned status %d", resp.StatusCode)
	}
	var out employeeV1.OnCallEmployeesResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		log.Errorf("employee.on_call decode_error err=%v", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return out.Employees, nil
}

func (c *clientImpl) CheckActiveEmergencies(ctx context.Context, employeeID uint) (bool, error) {
	log := c.logger.WithContext(ctx)
	endpoint := fmt.Sprintf("/api/v1/employees/%d/active-emergencies", employeeID)
	resp, err := c.retryGet(ctx, endpoint)
	if err != nil {
		log.Errorf("employee.active_emergencies http_error id=%d err=%v", employeeID, err)
		return false, fmt.Errorf("failed to check active emergencies for employee %d: %w", employeeID, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Errorf("employee.active_emergencies non_200 id=%d status=%d", employeeID, resp.StatusCode)
		return false, fmt.Errorf("employee service returned status %d", resp.StatusCode)
	}
	var out employeeV1.ActiveEmergenciesResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		log.Errorf("employee.active_emergencies decode_error id=%d err=%v", employeeID, err)
		return false, fmt.Errorf("failed to decode response: %w", err)
	}
	return out.HasActiveEmergencies, nil
}
