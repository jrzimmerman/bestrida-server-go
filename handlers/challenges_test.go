package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/mgo.v2/bson"

	log "github.com/Sirupsen/logrus"
	"github.com/jrzimmerman/bestrida-server-go/handlers"
	"github.com/jrzimmerman/bestrida-server-go/models"
)

func TestGetChallengeByIDSuccess(t *testing.T) {
	id := "57be4f7ef7fb96130084f0b2"

	// Create the http request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/challenges/%v", id), nil)
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
	var c models.Challenge
	if err := json.NewDecoder(rec.Body).Decode(&c); err != nil {
		t.Errorf("unable to decode response: %s", err)
	}

	log.WithField("Challenge ID", c.ID).Info("User returned from MongoDB")

	if c.ID != bson.ObjectIdHex(id) {
		t.Errorf("unexpected user")
	}
}

func TestGetChallengeByIDFailureID(t *testing.T) {
	id := "57fe7835bdb0181b8cfe0510"

	// Create the http request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/challenges/%v", id), nil)
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

func TestGetChallengeByIDFailureString(t *testing.T) {
	id := "bsonID"

	// Create the http request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/challenges/%v", id), nil)
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
