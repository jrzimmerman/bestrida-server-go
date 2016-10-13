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

func TestGetAthleteByIDFromStravaSuccess(t *testing.T) {
	id := 1027935

	// Create the http request
	req, err := http.NewRequest("GET", fmt.Sprintf("/strava/athletes/%v", id), nil)
	if err != nil {
		t.Fatal("unable to generate request", err)
	}

	// Send the request to the API
	rec := httptest.NewRecorder()
	handlers.API().ServeHTTP(rec, req)

	// Check the status code
	if exp := http.StatusOK; rec.Code != exp {
		t.Fatalf("expected status code %v, got: %v", exp, rec.Code)
	}

	// Unmarshal and check the response body
	var a *strava.AthleteDetailed
	if err := json.NewDecoder(rec.Body).Decode(&a); err != nil {
		t.Fatalf("unable to decode response: %s", err)
	}

	log.WithField("Athlete ID", a.Id).Info("Athlete returned from MongoDB")

	if a.Id != int64(id) {
		t.Fatalf("unexpected athlete")
	}

}
