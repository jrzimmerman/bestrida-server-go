package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/jrzimmerman/bestrida-server-go/handlers"
	strava "github.com/strava/go.strava"
)

// TestGetAthleteByIDFromStravaSuccess retrieves the athlete by ID from Strava
func TestGetAthleteByIDFromStravaSuccess(t *testing.T) {
	id := 1027935

	// Create the http request
	req, err := http.NewRequest("GET", fmt.Sprintf("/strava/athletes/%v", id), nil)
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
	var a *strava.AthleteDetailed
	if err := json.NewDecoder(rec.Body).Decode(&a); err != nil {
		t.Errorf("unable to decode response: %s", err)
	}

	log.WithField("Athlete ID", a.Id).Info("Athlete returned from Strava")

	if a.Id != int64(id) {
		t.Errorf("unexpected athlete")
	}
}

func TestGetAthleteByIDFromStravaFailureURL(t *testing.T) {
	id := "fred"

	// Create the http request
	req, err := http.NewRequest("GET", fmt.Sprintf("/strava/athletes/%v", id), nil)
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

func TestGetAthleteByIDFromStravaFailureID(t *testing.T) {
	id := 0

	// Create the http request
	req, err := http.NewRequest("GET", fmt.Sprintf("/strava/athletes/%v", id), nil)
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

func TestGetFriendsByUserIDFromStravaSuccess(t *testing.T) {
	id := 1027935

	// Create the http request
	req, err := http.NewRequest("GET", fmt.Sprintf("/strava/athletes/%v/friends", id), nil)
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
	var as *[]strava.AthleteSummary
	if err := json.NewDecoder(rec.Body).Decode(&as); err != nil {
		t.Errorf("unable to decode response: %s", err)
	}

	log.Info("Athlete friends returned from Strava")
}

func TestGetFriendsByUserIDFromStravaFailureURL(t *testing.T) {
	id := "fred"

	// Create the http request
	req, err := http.NewRequest("GET", fmt.Sprintf("/strava/athletes/%v/friends", id), nil)
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

func TestGetFriendsByUserIDFromStravaFailureID(t *testing.T) {
	id := 0

	// Create the http request
	req, err := http.NewRequest("GET", fmt.Sprintf("/strava/athletes/%v/friends", id), nil)
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

func TestGetSegmentsByUserIDFromStravaSuccess(t *testing.T) {
	id := 1027935

	// Create the http request
	req, err := http.NewRequest("GET", fmt.Sprintf("/strava/athletes/%v/segments", id), nil)
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
	var ss *[]strava.SegmentSummary
	if err := json.NewDecoder(rec.Body).Decode(&ss); err != nil {
		t.Errorf("unable to decode response: %s", err)
	}

	log.Info("Athlete segments returned from Strava")
}

// TestGetSegmentsByUserIDFromStravaFailureURL will test retrieving a user from strava with a bad athlete ID
func TestGetSegmentsByUserIDFromStravaFailureURL(t *testing.T) {
	id := "fred"

	// Create the http request
	req, err := http.NewRequest("GET", fmt.Sprintf("/strava/athletes/%v/segments", id), nil)
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

// TestGetSegmentsByUserIDFromStravaFailureID will test retrieving a user from strava with a bad athlete ID
func TestGetSegmentsByUserIDFromStravaFailureID(t *testing.T) {
	id := 0

	// Create the http request
	req, err := http.NewRequest("GET", fmt.Sprintf("/strava/athletes/%v/segments", id), nil)
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
