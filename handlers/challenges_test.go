package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/pressly/chi"
)

// TestGetChallengeByIDSuccess tests to successfully get a challenge by ID from the database
func TestGetChallengeByIDSuccess(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/:id", GetChallengeByID)
	server := httptest.NewServer(r)

	id := "57be4f7ef7fb96130084f0b2"

	// Create the http request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/"+id, server.URL), nil)
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
	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Errorf("unable to decode response: %s", err)
	}

	challenge := result.Content.(map[string]interface{})
	log.WithField("Challenge ID", challenge["id"]).Info("Challenge returned from DB")

	if challenge["id"] != id {
		t.Errorf("unexpected challenge")
	}
}

func TestGetChallengeByIDFailureID(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/:id", GetChallengeByID)
	server := httptest.NewServer(r)

	id := "57fe7835bdb0181b8cfe0510"

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

func TestGetChallengeByIDFailureString(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/:id", GetChallengeByID)
	server := httptest.NewServer(r)

	id := "bsonID"

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
