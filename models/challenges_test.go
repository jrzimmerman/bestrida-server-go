package models

import (
	"testing"

	"gopkg.in/mgo.v2/bson"
)

// TestGetChallengeByIDSuccess tests for a successfully found ObjectID
func TestGetChallengeByIDSuccess(t *testing.T) {
	id := bson.ObjectIdHex("59a2ed8cf0221020f76a5e9d")

	challenge, err := GetChallengeByID(id)
	if err != nil {
		t.Errorf("Unable to retrieve challenge by ObjectID:\n %v", err)
	}

	if challenge.ID != id {
		t.Errorf("Challenge ID %v, is not equal to %v", challenge.ID, id)
	}
}

// TestGetChallengeByIDFailure tests that an error is returned when an ObjectID is not found
func TestGetChallengeByIDFailure(t *testing.T) {
	id := bson.ObjectIdHex("000000000000000000000000")

	_, err := GetChallengeByID(id)
	if err.Error() != "not found" {
		t.Errorf("Unable to throw error for ID:\n %v", err)
	}
}

func TestCreateChallengeSuccess(t *testing.T) {
	id := bson.NewObjectId()
	c := Challenge{
		ID: id,
	}
	if err := CreateChallenge(c); err != nil {
		t.Fatalf("Error creating a new test challenge:\n %v", err)
	}
	defer RemoveChallenge(id)
}

func TestCreateChallengeFailure(t *testing.T) {
	c := Challenge{
		ID: "fred",
	}
	if err := CreateChallenge(c); err == nil {
		t.Error("Did not handle error creating a new test challenge")
	}
}
