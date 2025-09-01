package clients

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pd120424d/mountain-service/api/shared/auth"
	"github.com/pd120424d/mountain-service/api/shared/client"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/stretchr/testify/assert"
)

// fakeHTTPClient wraps the real HTTP client to capture URLs being called

type testHTTPClient struct{ inner *client.HTTPClient }

type errHTTPClient struct{}

func (t *testHTTPClient) Get(ctx context.Context, endpoint string) (*http.Response, error) {
	return t.inner.Get(ctx, endpoint)
}

func (e *errHTTPClient) Get(ctx context.Context, endpoint string) (*http.Response, error) {
	return nil, fmt.Errorf("boom")
}

func TestUrgencyClient_PathComposition(t *testing.T) {
	t.Parallel()
	log := utils.NewTestLogger()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{}"))
	}))
	defer server.Close()

	t.Run("it succeeds with base URL without /api/v1", func(t *testing.T) {
		base := server.URL
		httpClient := client.NewHTTPClient(client.HTTPClientConfig{BaseURL: base, Timeout: 2 * time.Second, ServiceAuth: auth.NewServiceAuth(auth.ServiceAuthConfig{Secret: "s", ServiceName: "activity-service", TokenTTL: time.Hour}), Logger: log})
		uc := &urgencyClient{httpClient: httpClient, logger: log}
		resp, err := uc.GetUrgencyByID(context.Background(), 2)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("it succeeds with base URL includes /api/v1", func(t *testing.T) {
		base := server.URL + "/api/v1"
		httpClient := client.NewHTTPClient(client.HTTPClientConfig{BaseURL: base, Timeout: 2 * time.Second, ServiceAuth: auth.NewServiceAuth(auth.ServiceAuthConfig{Secret: "s", ServiceName: "activity-service", TokenTTL: time.Hour}), Logger: log})
		uc := &urgencyClient{httpClient: httpClient, logger: log}
		resp, err := uc.GetUrgencyByID(context.Background(), 2)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("it returns error when urgency service returns non-200", func(t *testing.T) {
		server500 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server500.Close()
		httpClient := client.NewHTTPClient(client.HTTPClientConfig{BaseURL: server500.URL, Timeout: 2 * time.Second, ServiceAuth: auth.NewServiceAuth(auth.ServiceAuthConfig{Secret: "s", ServiceName: "activity-service", TokenTTL: time.Hour}), Logger: log})
		uc := &urgencyClient{httpClient: httpClient, logger: log}
		resp, err := uc.GetUrgencyByID(context.Background(), 42)
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
		uc := &urgencyClient{httpClient: httpClient, logger: log}
		resp, err := uc.GetUrgencyByID(context.Background(), 7)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("it returns error when HTTP client returns error", func(t *testing.T) {
		uc := &urgencyClient{httpClient: &errHTTPClient{}, logger: log}
		resp, err := uc.GetUrgencyByID(context.Background(), 1)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
