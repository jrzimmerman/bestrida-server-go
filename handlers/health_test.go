package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
)

func TestGetHealthCheckSuccess(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/", GetHealthCheck)
	server := httptest.NewServer(r)

	// Create the http request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/", server.URL), nil)
	if err != nil {
		t.Error("unable to generate request", err)
	}

	// Send the request to the API
	resp, err := http.DefaultClient.Do(req)

	// Check the status code
	if exp := http.StatusOK; resp.StatusCode != exp {
		t.Errorf("expected status code %v, got: %v", exp, resp.StatusCode)
	}
}
