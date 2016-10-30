package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/jrzimmerman/bestrida-server-go/handlers"
	"github.com/jrzimmerman/bestrida-server-go/models"
)

func TestGetUserByIDSuccess(t *testing.T) {
	id := 17198619

	// Create the http request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/users/%v", id), nil)
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

	// Unmarshal and check the response body
	var u models.User
	if err := json.NewDecoder(rec.Body).Decode(&u); err != nil {
		t.Errorf("unable to decode response: %s", err)
	}

	log.WithField("User ID", u.ID).Info("User returned from MongoDB")

	if u.ID != id {
		t.Errorf("unexpected user")
	}
}

func TestGetUserByIDFailureID(t *testing.T) {
	id := 0

	// Create the http request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/users/%v", id), nil)
	if err != nil {
		t.Error("unable to generate request", err)
	}

	// Send the request to the API
	rec := httptest.NewRecorder()
	handlers.API().ServeHTTP(rec, req)

	// Check the status code
	if exp := http.StatusInternalServerError; rec.Code != exp {
		t.Errorf("expected status code %v, got: %v", exp, rec.Code)
	}
}

func TestGetUserByIDFailureName(t *testing.T) {
	id := "fred"

	// Create the http request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/users/%v", id), nil)
	if err != nil {
		t.Error("unable to generate request", err)
	}

	// Send the request to the API
	rec := httptest.NewRecorder()
	handlers.API().ServeHTTP(rec, req)

	// Check the status code
	if exp := http.StatusInternalServerError; rec.Code != exp {
		t.Errorf("expected status code %v, got: %v", exp, rec.Code)
	}
}
