package models

import (
	"testing"

	"gopkg.in/mgo.v2/bson"
)

func TestGetChallengeByID(t *testing.T) {
	id := bson.ObjectIdHex("57be4f7ef7fb96130084f0b2")

	challenge, err := GetChallengeByID(id)
	if err != nil {
		t.Errorf("Unable to retrieve challenge by ObjectID:\n %v", err)
	}

	if challenge.ID != id {
		t.Errorf("Challenge ID %v, is not equal to %v", challenge.ID, id)
	}
}
