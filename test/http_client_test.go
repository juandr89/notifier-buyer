package service_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/juandr89/delivery-notifier-buyer/server"
	"github.com/stretchr/testify/assert"
)

func TestDoRequestWithRetry(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`OK`))
	}))
	opts := server.RequestOptions{
		Method:         "GET",
		URL:            mockServer.URL,
		Headers:        map[string]string{"Content-Type": "application/json"},
		Body:           nil,
		RequestTimeout: 2 * time.Second,
		MaxRetries:     3,
		RetryDelay:     1 * time.Second,
	}
	t.Run("Success", func(t *testing.T) {
		defer mockServer.Close()

		resp, err := server.DoRequestWithRetry(opts)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("AllRetriesFail", func(t *testing.T) {
		defer mockServer.Close()

		resp, err := server.DoRequestWithRetry(opts)

		assert.Nil(t, resp)
		assert.Error(t, err)

		assert.Contains(t, err.Error(), "request failed after 3 attempts")
	})

}

func TestNewRouter(t *testing.T) {
	port := ":8080"

	// Call the NewRouter function to create a new server
	serverInstance := server.NewRouter(port)

	// Assertions to verify the server instance
	assert.NotNil(t, serverInstance)
	assert.Equal(t, port, serverInstance.Addr)

	// Check that the handler is of type *mux.Router
	router, ok := serverInstance.Handler.(*mux.Router)
	assert.True(t, ok, "Handler should be of type *mux.Router")

	// Additional checks can be added here if necessary (e.g., routes)
	assert.NotNil(t, router, "Router should not be nil")
}
