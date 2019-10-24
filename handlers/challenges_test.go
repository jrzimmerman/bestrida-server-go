package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
)

// TestGetChallengeByIDSuccess tests to successfully get a challenge by ID from the database
// func TestGetChallengeByIDSuccess(t *testing.T) {
// 	r := chi.NewRouter()
// 	r.Get("/{id}", GetChallengeByID)
// 	server := httptest.NewServer(r)

// 	id := "59a309ddf02210361b3b027f"

// 	// Create the http request
// 	req, err := http.NewRequest("GET", fmt.Sprintf("%s/"+id, server.URL), nil)
// 	if err != nil {
// 		t.Error("unable to generate request", err)
// 	}

// 	// Send the request to the API
// 	resp, err := http.DefaultClient.Do(req)

// 	// Check the status code
// 	if exp := http.StatusOK; resp.StatusCode != exp {
// 		t.Errorf("expected status code %v, got: %v", exp, resp.StatusCode)
// 	}

// 	// Unmarshal and check the response body
// 	var c models.Challenge
// 	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
// 		t.Errorf("unable to decode response: %s", err)
// 	}

// 	log.WithField("Challenge ID", c.ID).Info("User returned from database")

// 	if c.ID != bson.ObjectIdHex(id) {
// 		t.Errorf("unexpected user")
// 	}
// }

func TestGetChallengeByIDFailureID(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/{id}", GetChallengeByID)
	server := httptest.NewServer(r)

	id := "000000000000000000000000"

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
	r.Get("/{id}", GetChallengeByID)
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

// func TestCreateChallengeSuccess(t *testing.T) {
// 	r := chi.NewRouter()
// 	r.Post("/", CreateChallenge)
// 	server := httptest.NewServer(r)

// 	var segmentID int64 = 12924664
// 	byteReq := []byte(`
// 		{
// 			"challengerId": 17198619,
// 			"challengeeId": 1027935,
// 			"segmentId": 12924664,
// 			"creationDate": "2016-08-28T02:38:20.926Z",
// 			"completionDate": "2017-08-27T02:38:20.926Z"
// 		}
// 	`)
// 	log.Infof("byteReq:\n%v", byteReq)
// 	body := bytes.NewReader(byteReq)

// 	// Create the http request
// 	req, err := http.NewRequest("POST", fmt.Sprintf("%s/", server.URL), body)
// 	if err != nil {
// 		t.Error("unable to generate request", err)
// 	}

// 	// Send the request to the API
// 	resp, err := http.DefaultClient.Do(req)

// 	// Check the status code
// 	if exp := http.StatusOK; resp.StatusCode != exp {
// 		t.Errorf("expected status code %v, got: %v", exp, resp.StatusCode)
// 	}

// 	// Unmarshal and check the response body
// 	var c models.Challenge
// 	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
// 		t.Errorf("unable to decode response: %s", err)
// 	}

// 	log.WithField("Challenge ID", c.ID).Info("Challenge created")

// 	if c.Segment.ID != segmentID {
// 		t.Errorf("unexpected segment id from create challenge")
// 	}

// 	if err := models.RemoveChallenge(c.ID); err != nil {
// 		t.Errorf("unable to remove challenge: %s", err)
// 	}
// }
