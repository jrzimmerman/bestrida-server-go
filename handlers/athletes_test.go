package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	"github.com/strava/go.strava"
)

// TestGetAthleteByIDFromStravaSuccess retrieves the athlete by ID from Strava
func TestGetAthleteByIDFromStravaSuccess(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/{id}", GetAthleteByIDFromStrava)
	server := httptest.NewServer(r)

	id := 17198619

	// Create the http request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/"+strconv.Itoa(id), server.URL), nil)
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
	var a *strava.AthleteDetailed
	if err := json.NewDecoder(resp.Body).Decode(&a); err != nil {
		t.Errorf("unable to decode response: %s", err)
	}

	log.WithField("Athlete ID", a.Id).Info("Athlete returned from Strava")

	if a.Id != int64(id) {
		t.Errorf("unexpected athlete")
	}
}

func TestGetAthleteByIDFromStravaFailureURL(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/{id}", GetAthleteByIDFromStrava)
	server := httptest.NewServer(r)

	id := "fred"

	// Create the http request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/"+id, server.URL), nil)
	if err != nil {
		t.Error("unable to generate request", err)
	}

	// Send the request to the API
	resp, err := http.DefaultClient.Do(req)

	// Check the status code
	if exp := http.StatusBadRequest; resp.StatusCode != exp {
		t.Errorf("expected status code %v, got: %v", exp, resp.StatusCode)
	}
}

// func TestGetAthleteByIDFromStravaFailureID(t *testing.T) {
// 	id := 0

// 	// Create the http request
// 	req, err := http.NewRequest("GET", fmt.Sprintf("/strava/athletes/%v", id), nil)
// 	if err != nil {
// 		t.Error("unable to generate request", err)
// 	}

// 	// Send the request to the API
// 	rec := httptest.NewRecorder()
// 	API().ServeHTTP(rec, req)

// 	// Check the status code
// 	if exp := http.StatusInternalServerError; rec.Code != exp {
// 		t.Errorf("expected status code %v, got: %v", exp, rec.Code)
// 	}
// }

// func TestGetFriendsByUserIDFromStravaSuccess(t *testing.T) {
// 	id := 17198619

// 	// Create the http request
// 	req, err := http.NewRequest("GET", fmt.Sprintf("/strava/athletes/%v/friends", id), nil)
// 	if err != nil {
// 		t.Error("unable to generate request", err)
// 	}

// 	// Send the request to the API
// 	rec := httptest.NewRecorder()
// 	API().ServeHTTP(rec, req)

// 	// Check the status code
// 	if exp := http.StatusOK; rec.Code != exp {
// 		t.Errorf("expected status code %v, got: %v", exp, rec.Code)
// 	}

// 	// Unmarshal and check the response body
// 	var as *[]strava.AthleteSummary
// 	if err := json.NewDecoder(rec.Body).Decode(&as); err != nil {
// 		t.Errorf("unable to decode response: %s", err)
// 	}

// 	log.Info("Athlete friends returned from Strava")
// }

func TestGetFriendsByIDFromStravaFailureURL(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/{id}", GetFriendsByUserIDFromStrava)
	server := httptest.NewServer(r)

	id := "fred"

	// Create the http request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/"+id, server.URL), nil)
	if err != nil {
		t.Error("unable to generate request", err)
	}

	// Send the request to the API
	resp, err := http.DefaultClient.Do(req)

	// Check the status code
	if exp := http.StatusBadRequest; resp.StatusCode != exp {
		t.Errorf("expected status code %v, got: %v", exp, resp.StatusCode)
	}
}

// func TestGetFriendsByUserIDFromStravaFailureID(t *testing.T) {
// 	id := 0

// 	// Create the http request
// 	req, err := http.NewRequest("GET", fmt.Sprintf("/strava/athletes/%v/friends", id), nil)
// 	if err != nil {
// 		t.Error("unable to generate request", err)
// 	}

// 	// Send the request to the API
// 	rec := httptest.NewRecorder()
// 	API().ServeHTTP(rec, req)

// 	// Check the status code
// 	if exp := http.StatusInternalServerError; rec.Code != exp {
// 		t.Errorf("expected status code %v, got: %v", exp, rec.Code)
// 	}
// }

func TestGetSegmentsByUserIDFromStravaSuccess(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/{id}", GetSegmentsByUserIDFromStrava)
	server := httptest.NewServer(r)

	id := 17198619

	// Create the http request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/"+strconv.Itoa(id), server.URL), nil)
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
	var ss *[]strava.SegmentSummary
	if err := json.NewDecoder(resp.Body).Decode(&ss); err != nil {
		t.Errorf("unable to decode response: %s", err)
	}

	log.Info("Athlete segments returned from Strava")
}

// // TestGetSegmentsByUserIDFromStravaFailureURL will test retrieving a user from strava with a bad athlete ID
// func TestGetSegmentsByUserIDFromStravaFailureURL(t *testing.T) {
// 	id := "fred"

// 	// Create the http request
// 	req, err := http.NewRequest("GET", fmt.Sprintf("/strava/athletes/%v/segments", id), nil)
// 	if err != nil {
// 		t.Error("unable to generate request", err)
// 	}

// 	// Send the request to the API
// 	rec := httptest.NewRecorder()
// 	API().ServeHTTP(rec, req)

// 	// Check the status code
// 	if exp := http.StatusInternalServerError; rec.Code != exp {
// 		t.Errorf("expected status code %v, got: %v", exp, rec.Code)
// 	}
// }

// // TestGetSegmentsByUserIDFromStravaFailureID will test retrieving a user from strava with a bad athlete ID
// func TestGetSegmentsByUserIDFromStravaFailureID(t *testing.T) {
// 	id := 0

// 	// Create the http request
// 	req, err := http.NewRequest("GET", fmt.Sprintf("/strava/athletes/%v/segments", id), nil)
// 	if err != nil {
// 		t.Error("unable to generate request", err)
// 	}

// 	// Send the request to the API
// 	rec := httptest.NewRecorder()
// 	API().ServeHTTP(rec, req)

// 	// Check the status code
// 	if exp := http.StatusInternalServerError; rec.Code != exp {
// 		t.Errorf("expected status code %v, got: %v", exp, rec.Code)
// 	}
// }
