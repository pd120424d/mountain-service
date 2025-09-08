package clients

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pd120424d/mountain-service/api/shared/auth"
	"github.com/pd120424d/mountain-service/api/shared/client"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/stretchr/testify/assert"
)

func TestEmployeeClient_GetEmployeeByID(t *testing.T) {
	t.Parallel()
	log := utils.NewTestLogger()

	// Shared 200 server (no path assertion to mirror urgency tests and avoid double /api/v1 coupling)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{}"))
	}))
	defer server.Close()

	t.Run("it succeeds with base URL without /api/v1", func(t *testing.T) {
		base := server.URL
		httpClient := client.NewHTTPClient(client.HTTPClientConfig{BaseURL: base, Timeout: 2 * time.Second, ServiceAuth: auth.NewServiceAuth(auth.ServiceAuthConfig{Secret: "s", ServiceName: "activity-service", TokenTTL: time.Hour}), Logger: log})
		ec := &employeeClient{httpClient: httpClient, logger: log}
		resp, err := ec.GetEmployeeByID(t.Context(), 2)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("it succeeds with base URL includes /api/v1", func(t *testing.T) {
		base := server.URL + "/api/v1"
		httpClient := client.NewHTTPClient(client.HTTPClientConfig{BaseURL: base, Timeout: 2 * time.Second, ServiceAuth: auth.NewServiceAuth(auth.ServiceAuthConfig{Secret: "s", ServiceName: "activity-service", TokenTTL: time.Hour}), Logger: log})
		ec := &employeeClient{httpClient: httpClient, logger: log}
		resp, err := ec.GetEmployeeByID(t.Context(), 2)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("it returns error when employee service returns non-200", func(t *testing.T) {
		server500 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server500.Close()
		httpClient := client.NewHTTPClient(client.HTTPClientConfig{BaseURL: server500.URL, Timeout: 2 * time.Second, ServiceAuth: auth.NewServiceAuth(auth.ServiceAuthConfig{Secret: "s", ServiceName: "activity-service", TokenTTL: time.Hour}), Logger: log})
		ec := &employeeClient{httpClient: httpClient, logger: log}
		resp, err := ec.GetEmployeeByID(t.Context(), 42)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("it returns error when response body is invalid JSON", func(t *testing.T) {
		serverInvalid := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte("not-json"))
		}))
		defer serverInvalid.Close()
		httpClient := client.NewHTTPClient(client.HTTPClientConfig{BaseURL: serverInvalid.URL, Timeout: 2 * time.Second, ServiceAuth: auth.NewServiceAuth(auth.ServiceAuthConfig{Secret: "s", ServiceName: "activity-service", TokenTTL: time.Hour}), Logger: log})
		ec := &employeeClient{httpClient: httpClient, logger: log}
		resp, err := ec.GetEmployeeByID(t.Context(), 7)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("it returns error when HTTP client returns error", func(t *testing.T) {
		ec := &employeeClient{httpClient: &errHTTPClient{}, logger: log}
		resp, err := ec.GetEmployeeByID(t.Context(), 1)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("it times out via HTTP client timeout", func(t *testing.T) {
		serverSlow := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(200 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("{}"))
		}))
		defer serverSlow.Close()
		httpClient := client.NewHTTPClient(client.HTTPClientConfig{BaseURL: serverSlow.URL, Timeout: 50 * time.Millisecond, ServiceAuth: auth.NewServiceAuth(auth.ServiceAuthConfig{Secret: "s", ServiceName: "activity-service", TokenTTL: time.Hour}), Logger: log})
		ec := &employeeClient{httpClient: httpClient, logger: log}
		resp, err := ec.GetEmployeeByID(t.Context(), 3)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("it respects context cancellation deadline", func(t *testing.T) {
		serverBlock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(200 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("{}"))
		}))
		defer serverBlock.Close()
		httpClient := client.NewHTTPClient(client.HTTPClientConfig{BaseURL: serverBlock.URL, Timeout: 2 * time.Second, ServiceAuth: auth.NewServiceAuth(auth.ServiceAuthConfig{Secret: "s", ServiceName: "activity-service", TokenTTL: time.Hour}), Logger: log})
		ec := &employeeClient{httpClient: httpClient, logger: log}
		ctx, cancel := context.WithTimeout(t.Context(), 30*time.Millisecond)
		defer cancel()
		resp, err := ec.GetEmployeeByID(ctx, 4)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
