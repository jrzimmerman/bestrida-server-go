package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/jrzimmerman/bestrida-server-go/models"
)

func TestGetUserByID(t *testing.T) {
	id := 1027935

	// Create the http request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/users/%v", id), nil)
	if err != nil {
		t.Fatal("unable to generate request", err)
	}

	// Send the request to the API
	rec := httptest.NewRecorder()
	API().ServeHTTP(rec, req)

	// Check the status code
	if exp := http.StatusOK; rec.Code != exp {
		t.Fatalf("expected status code %v, got: %v", exp, rec.Code)
	}

	// Unmarshal and check the response body
	var u models.User
	if err := json.NewDecoder(rec.Body).Decode(&u); err != nil {
		t.Fatalf("unable to decode response: %s", err)
	}

	log.WithField("User ID", u.ID).Info("User returned from MongoDB")

	if u.ID != id {
		t.Fatalf("unexpected user")
	}
}
