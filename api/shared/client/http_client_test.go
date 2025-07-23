package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pd120424d/mountain-service/api/shared/auth"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_Get(t *testing.T) {
	t.Parallel()

	log := utils.NewTestLogger()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	serviceAuth := auth.NewServiceAuth(auth.ServiceAuthConfig{
		Secret:      "test-secret",
		ServiceName: "test-service",
		TokenTTL:    1 * time.Hour,
	})
	httpClient := NewHTTPClient(HTTPClientConfig{
		BaseURL:     server.URL,
		Timeout:     30 * time.Second,
		ServiceAuth: serviceAuth,
		Logger:      log,
	})
	defer server.Close()

	t.Run("it succeeds with a GET request", func(t *testing.T) {

		resp, err := httpClient.Get(context.Background(), "/test")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("it suceeds with a POST request", func(t *testing.T) {
		resp, err := httpClient.Post(context.Background(), "/test", nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("it suceeds with a PUT request", func(t *testing.T) {
		resp, err := httpClient.Put(context.Background(), "/test", nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("it suceeds with a DELETE request", func(t *testing.T) {
		resp, err := httpClient.Delete(context.Background(), "/test")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
