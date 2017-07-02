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

func TestGetUserByIDSuccess(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/:id", GetUserByID)
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
	var u models.User
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		t.Errorf("unable to decode response: %s", err)
	}

	log.WithField("User ID", u.ID).Info("User returned from database")

	if u.ID != int64(id) {
		t.Errorf("unexpected user")
	}
}

func TestGetUserByIDFailureID(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/:id", GetUserByID)
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

func TestGetUserByIDFailureName(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/:id", GetUserByID)
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
	if exp := http.StatusInternalServerError; resp.StatusCode != exp {
		t.Errorf("expected status code %v, got: %v", exp, resp.StatusCode)
	}
}
