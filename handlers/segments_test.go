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

func TestGetSegmentByID(t *testing.T) {
	id := 2539276

	// Create the http request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/segments/%v", id), nil)
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
	var s models.Segment
	if err := json.NewDecoder(rec.Body).Decode(&s); err != nil {
		t.Fatalf("unable to decode response: %s", err)
	}

	log.WithField("Segment ID", s.ID).Info("Segment returned from MongoDB")

	if s.ID != id {
		t.Fatalf("unexpected segment")
	}
}
