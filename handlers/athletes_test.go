package handlers

// // TestGetAthleteByIDFromStravaSuccess retrieves the athlete by ID from Strava
// func TestGetAthleteByIDFromStravaSuccess(t *testing.T) {
// 	r := chi.NewRouter()
// 	r.Get("/:id", GetAthleteByIDFromStrava)
// 	server := httptest.NewServer(r)

// 	id := 1027935

// 	// Create the http request
// 	req, err := http.NewRequest("GET", fmt.Sprintf("%s/"+strconv.Itoa(id), server.URL), nil)
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
// 	var result Response
// 	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
// 		t.Errorf("unable to decode response: %s", err)
// 	}

// 	athlete := result.Content.(map[string]interface{})
// 	log.WithField("Athlete ID", athlete["id"]).Info("Athlete successfully returned from Strava")

// 	if athlete["id"] != id {
// 		t.Errorf("unexpected athlete")
// 	}
// }

// func TestGetAthleteByIDFromStravaFailureURL(t *testing.T) {
// 	id := "fred"

// 	// Create the http request
// 	req, err := http.NewRequest("GET", fmt.Sprintf("/strava/athletes/%v", id), nil)
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

// func TestGetAthleteByIDFromStravaFailureID(t *testing.T) {
// 	id := 0

// 	// Create the http request
// 	req, err := http.NewRequest("GET", fmt.Sprintf("/strava/athletes/%v", id), nil)
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

// func TestGetFriendsByUserIDFromStravaSuccess(t *testing.T) {
// 	id := 1027935

// 	// Create the http request
// 	req, err := http.NewRequest("GET", fmt.Sprintf("/strava/athletes/%v/friends", id), nil)
// 	if err != nil {
// 		t.Error("unable to generate request", err)
// 	}

// 	// Send the request to the API
// 	rec := httptest.NewRecorder()
// 	handlers.API().ServeHTTP(rec, req)

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

// func TestGetFriendsByUserIDFromStravaFailureURL(t *testing.T) {
// 	id := "fred"

// 	// Create the http request
// 	req, err := http.NewRequest("GET", fmt.Sprintf("/strava/athletes/%v/friends", id), nil)
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

// func TestGetFriendsByUserIDFromStravaFailureID(t *testing.T) {
// 	id := 0

// 	// Create the http request
// 	req, err := http.NewRequest("GET", fmt.Sprintf("/strava/athletes/%v/friends", id), nil)
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

// func TestGetSegmentsByUserIDFromStravaSuccess(t *testing.T) {
// 	id := 1027935

// 	// Create the http request
// 	req, err := http.NewRequest("GET", fmt.Sprintf("/strava/athletes/%v/segments", id), nil)
// 	if err != nil {
// 		t.Error("unable to generate request", err)
// 	}

// 	// Send the request to the API
// 	rec := httptest.NewRecorder()
// 	handlers.API().ServeHTTP(rec, req)

// 	// Check the status code
// 	if exp := http.StatusOK; rec.Code != exp {
// 		t.Errorf("expected status code %v, got: %v", exp, rec.Code)
// 	}

// 	// Unmarshal and check the response body
// 	var ss *[]strava.SegmentSummary
// 	if err := json.NewDecoder(rec.Body).Decode(&ss); err != nil {
// 		t.Errorf("unable to decode response: %s", err)
// 	}

// 	log.Info("Athlete segments returned from Strava")
// }

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
// 	handlers.API().ServeHTTP(rec, req)

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
// 	handlers.API().ServeHTTP(rec, req)

// 	// Check the status code
// 	if exp := http.StatusInternalServerError; rec.Code != exp {
// 		t.Errorf("expected status code %v, got: %v", exp, rec.Code)
// 	}
// }
