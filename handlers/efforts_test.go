package handlers

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/jrzimmerman/bestrida-server-go/handlers"
// 	strava "github.com/strava/go.strava"
// )

// func TestGetEffortsBySegmentIDFromStravaSuccess(t *testing.T) {
// 	id := 1027935
// 	segmentID := 9719730

// 	// Create the http request
// 	req, err := http.NewRequest("GET", fmt.Sprintf("/strava/athletes/%v/segments/%v/efforts", id, segmentID), nil)
// 	if err != nil {
// 		t.Error("unable to generate request", err)
// 	}

// 	// Send the request to the API
// 	rec := httptest.NewRecorder()
// 	handlers.API().ServeHTTP(rec, req)

// 	// Check the status code
// 	if exp := http.StatusOK; rec.Code != exp {
// 		t.Errorf("expected status code %v, got: %v", exp, rec.Code)
// 	}

// 	// Unmarshal and check the response body
// 	var efforts []*strava.EffortSummary
// 	if err := json.NewDecoder(rec.Body).Decode(&efforts); err != nil {
// 		t.Errorf("unable to decode response: %s", err)
// 	}

// 	if len(efforts) == 0 {
// 		t.Errorf("no efforts returned")
// 	}
// }

// func TestGetEffortsBySegmentIDFromStravaFailureUserID(t *testing.T) {
// 	id := 0
// 	segmentID := 9719730

// 	// Create the http request
// 	req, err := http.NewRequest("GET", fmt.Sprintf("/strava/athletes/%v/segments/%v/efforts", id, segmentID), nil)
// 	if err != nil {
// 		t.Error("unable to generate request", err)
// 	}

// 	// Send the request to the API
// 	rec := httptest.NewRecorder()
// 	handlers.API().ServeHTTP(rec, req)

// 	// Check the status code
// 	if exp := http.StatusInternalServerError; rec.Code != exp {
// 		t.Errorf("expected status code %v, got: %v", exp, rec.Code)
// 	}
// }

// func TestGetEffortsBySegmentIDFromStravaFailureUserInput(t *testing.T) {
// 	id := "test"
// 	segmentID := 9719730

// 	// Create the http request
// 	req, err := http.NewRequest("GET", fmt.Sprintf("/strava/athletes/%v/segments/%v/efforts", id, segmentID), nil)
// 	if err != nil {
// 		t.Error("unable to generate request", err)
// 	}

// 	// Send the request to the API
// 	rec := httptest.NewRecorder()
// 	handlers.API().ServeHTTP(rec, req)

// 	// Check the status code
// 	if exp := http.StatusInternalServerError; rec.Code != exp {
// 		t.Errorf("expected status code %v, got: %v", exp, rec.Code)
// 	}
// }

// func TestGetEffortsBySegmentIDFromStravaFailureSegmentInput(t *testing.T) {
// 	id := 1027935
// 	segmentID := "test"

// 	// Create the http request
// 	req, err := http.NewRequest("GET", fmt.Sprintf("/strava/athletes/%v/segments/%v/efforts", id, segmentID), nil)
// 	if err != nil {
// 		t.Error("unable to generate request", err)
// 	}

// 	// Send the request to the API
// 	rec := httptest.NewRecorder()
// 	handlers.API().ServeHTTP(rec, req)

// 	// Check the status code
// 	if exp := http.StatusInternalServerError; rec.Code != exp {
// 		t.Errorf("expected status code %v, got: %v", exp, rec.Code)
// 	}
// }
