package models

import "testing"

func TestGetUserByID(t *testing.T) {
	id := 1027935

	user, err := GetUserByID(id)
	if err != nil {
		t.Errorf("Unable to retrieve user by ID:\n %v", err)
	}

	if user.ID != id {
		t.Errorf("User ID %v, is not equal to %v", user.ID, id)
	}
}
