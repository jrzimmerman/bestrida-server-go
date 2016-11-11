package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jrzimmerman/bestrida-server-go/handlers"
)

func TestStravaAuth(t *testing.T) {
	// Create the http request
	req, err := http.NewRequest("GET", fmt.Sprintf("/strava/auth"), nil)
	if err != nil {
		t.Error("unable to generate request", err)
	}

	// Send the request to the API
	rec := httptest.NewRecorder()
	handlers.API().ServeHTTP(rec, req)

	// Check the status code
	if exp := http.StatusOK; rec.Code != exp {
		t.Errorf("expected status code %v, got: %v", exp, rec.Code)
	}
}
