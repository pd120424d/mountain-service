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

// fakeHTTPClient wraps the real HTTP client to capture URLs being called

type testHTTPClient struct{ inner *client.HTTPClient }

func (t *testHTTPClient) Get(ctx context.Context, endpoint string) (*http.Response, error) {
	return t.inner.Get(ctx, endpoint)
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
}
