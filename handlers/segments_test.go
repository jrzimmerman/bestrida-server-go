package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi"
	"github.com/jrzimmerman/bestrida-server-go/models"
	log "github.com/sirupsen/logrus"
)

func TestGetSegmentByIDSuccess(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/{id}", GetSegmentByID)
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

	log.WithField("Segment ID", s.ID).Info("Segment returned from database")

	if s.ID != int64(id) {
		t.Errorf("unexpected segment")
	}
}

func TestGetSegmentByIDFailureID(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/{id}", GetSegmentByID)
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
	r.Get("/{id}", GetSegmentByID)
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

func TestGetSegmentByIDFromStravaSuccess(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/{id}", GetSegmentByIDFromStrava)
	server := httptest.NewServer(r)

	// Hawk Hill segment ID
	id := 229781

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

	log.WithField("Segment ID", s.ID).Info("Segment returned from database")

	if s.ID != int64(id) {
		t.Errorf("unexpected segment")
	}
}

func TestGetSegmentByIDFromStravaFailureID(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/{id}", GetSegmentByIDFromStrava)
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

func TestGetSegmentByIDFromStravaFailureInput(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/{id}", GetSegmentByIDFromStrava)
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
