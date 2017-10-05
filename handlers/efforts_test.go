package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/strava/go.strava"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

func TestGetEffortsBySegmentIDFromStravaWithUserIDSuccess(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/{id}/{segmentID}", GetEffortsBySegmentIDFromStravaWithUserID)
	server := httptest.NewServer(r)

	// User ID
	id := 1027935
	segmentID := 12967163

	// Create the http request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/"+strconv.Itoa(id)+"/"+strconv.Itoa(segmentID), server.URL), nil)
	if err != nil {
		t.Error("unable to generate request", err)
	}

	// Send the request to the API
	resp, err := http.DefaultClient.Do(req)

	// Check the status code
	if exp := http.StatusOK; resp.StatusCode != exp {
		t.Errorf("expected status code %v, got: %v", exp, resp.StatusCode)
	}

	// Unmarshal and check the response body
	var efforts []*strava.SegmentEffortSummary
	if err := json.NewDecoder(resp.Body).Decode(&efforts); err != nil {
		t.Errorf("unable to decode response: %s", err)
	}

	log.Info("Segment efforts returned from Strava")

	if len(efforts) <= 0 {
		t.Errorf("segment efforts not found")
	}
}

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
// 	id := 17198619
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
