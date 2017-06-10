package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/jrzimmerman/bestrida-server-go/models"
	"github.com/pressly/chi"
)

func TestGetSegmentByIDSuccess(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/:id", GetSegmentByID)
	server := httptest.NewServer(r)

	id := 2539276

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
	var s models.Segment
	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		t.Errorf("unable to decode response: %s", err)
	}

	log.WithField("Segment ID", s.ID).Info("Segment returned from MongoDB")

	if s.ID != id {
		t.Errorf("unexpected segment")
	}
}

func TestGetSegmentByIDFailureID(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/:id", GetSegmentByID)
	server := httptest.NewServer(r)

	id := 0

	// Create the http request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/"+strconv.Itoa(id), server.URL), nil)
	if err != nil {
		t.Error("unable to generate request", err)
	}

	// Send the request to the API
	resp, err := http.DefaultClient.Do(req)

	// Check the status code
	if exp := http.StatusInternalServerError; resp.StatusCode != exp {
		t.Errorf("expected status code %v, got: %v", exp, resp.StatusCode)
	}
}

func TestGetSegmentByIDFailureInput(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/:id", GetSegmentByID)
	server := httptest.NewServer(r)

	id := "test"

	// Create the http request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/"+id, server.URL), nil)
	if err != nil {
		t.Error("unable to generate request", err)
	}

	// Send the request to the API
	resp, err := http.DefaultClient.Do(req)

	// Check the status code
	if exp := http.StatusInternalServerError; resp.StatusCode != exp {
		t.Errorf("expected status code %v, got: %v", exp, resp.StatusCode)
	}
}

// func TestGetSegmentByIDFromStravaSuccess(t *testing.T) {
// 	// Hawk Hill segment ID
// 	var id int64 = 229781

// 	// Create the http request
// 	req, err := http.NewRequest("GET", fmt.Sprintf("/strava/segments/%v", id), nil)
// 	if err != nil {
// 		t.Error("unable to generate request", err)
// 	}

// 	// Send the request to the API
// 	rec := httptest.NewRecorder()
// 	// handlers.API().ServeHTTP(rec, req)

// 	// Check the status code
// 	if exp := http.StatusOK; rec.Code != exp {
// 		t.Errorf("expected status code %v, got: %v", exp, rec.Code)
// 	}

// 	// Unmarshal and check the response body
// 	var s strava.SegmentDetailed
// 	if err := json.NewDecoder(rec.Body).Decode(&s); err != nil {
// 		t.Errorf("unable to decode response: %s", err)
// 	}

// 	log.WithField("Segment ID", s.Id).Info("Segment returned from MongoDB")

// 	if s.Id != id {
// 		t.Errorf("Expected segment ID %v, got %v instead", id, s.Id)
// 	}
// }

// func TestGetSegmentByIDFromStravaFailureID(t *testing.T) {
// 	id := 0

// 	// Create the http request
// 	req, err := http.NewRequest("GET", fmt.Sprintf("/strava/segments/%v", id), nil)
// 	if err != nil {
// 		t.Error("unable to generate request", err)
// 	}

// 	// Send the request to the API
// 	rec := httptest.NewRecorder()
// 	// handlers.API().ServeHTTP(rec, req)

// 	// Check the status code
// 	if exp := http.StatusInternalServerError; rec.Code != exp {
// 		t.Errorf("expected status code %v, got: %v", exp, rec.Code)
// 	}
// }

// func TestGetSegmentByIDFromStravaFailureInput(t *testing.T) {
// 	id := "test"

// 	// Create the http request
// 	req, err := http.NewRequest("GET", fmt.Sprintf("/strava/segments/%v", id), nil)
// 	if err != nil {
// 		t.Error("unable to generate request", err)
// 	}

// 	// Send the request to the API
// 	rec := httptest.NewRecorder()
// 	// handlers.API().ServeHTTP(rec, req)

// 	// Check the status code
// 	if exp := http.StatusInternalServerError; rec.Code != exp {
// 		t.Errorf("expected status code %v, got: %v", exp, rec.Code)
// 	}
// }
