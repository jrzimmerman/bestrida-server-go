package models

import (
	"testing"

	"gopkg.in/mgo.v2/bson"
)

func TestGetChallengeByIDSuccess(t *testing.T) {
	id := bson.ObjectIdHex("57be4f7ef7fb96130084f0b2")

	challenge, err := GetChallengeByID(id)
	if err != nil {
		t.Errorf("Unable to retrieve challenge by ObjectID:\n %v", err)
	}

	if challenge.ID != id {
		t.Errorf("Challenge ID %v, is not equal to %v", challenge.ID, id)
	}
}

func TestGetChallengeByIDFailure(t *testing.T) {
	id := bson.ObjectIdHex("000000000000000000000000")

	_, err := GetChallengeByID(id)
	if err.Error() != "not found" {
		t.Errorf("Unable to throw error for ID:\n %v", err)
	}
}
