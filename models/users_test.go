package models

import "testing"

func TestGetUserByIDSuccess(t *testing.T) {
	id := 1027935

	user, err := GetUserByID(id)
	if err != nil {
		t.Errorf("Unable to retrieve user by ID:\n %v", err)
	}

	if user.ID != id {
		t.Errorf("User ID %v, is not equal to %v", user.ID, id)
	}
}

func TestGetUserByIDFailure(t *testing.T) {
	id := 0

	_, err := GetUserByID(id)
	if err.Error() != "not found" {
		t.Errorf("Unable to throw error for ID:\n %v", err)
	}
}

func TestModifySegmentCount(t *testing.T) {
	// id for specific user
	id := 1027935

	// segmentID for a segment's strava segment ID
	segmentID := 10599051

	user, err := GetUserByID(id)
	if err != nil {
		t.Errorf("Unable to retrieve user by ID:\n %v", err)
	}

	var userSegmentCount int
	for i := range user.Segments {
		if user.Segments[i].ID == segmentID {
			userSegmentCount = user.Segments[i].Count
			break
		}
	}

	// modify the segment count for a users segment array by 1
	user.ModifySegmentCount(segmentID, 1)

	var modifiedSegmentCount int
	for i := range user.Segments {
		if user.Segments[i].ID == segmentID {
			modifiedSegmentCount = user.Segments[i].Count
			break
		}
	}

	if modifiedSegmentCount != userSegmentCount+1 {
		t.Errorf("Expected Modified User Segments to be %v, but received %v instead", userSegmentCount, modifiedSegmentCount)
	}

	// clean up users segment count after test
	user.ModifySegmentCount(segmentID, -1)

	var cleanedSegmentCount int
	for i := range user.Segments {
		if user.Segments[i].ID == segmentID {
			cleanedSegmentCount = user.Segments[i].Count
			break
		}
	}

	if cleanedSegmentCount != userSegmentCount {
		t.Errorf("Expected Cleaned User Segments to be %v, but received %v instead", userSegmentCount, cleanedSegmentCount)
	}
}
