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
