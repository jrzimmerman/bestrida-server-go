package models

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// Friend struct handles the MongoDB schema for each users friends
type Friend struct {
	ID             int    `json:"id"`
	UserName       string `json:"userName"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	FullName       string `json:"fullName"`
	Photo          string `json:"photo"`
	ChallengeCount int    `json:"challengeCount"`
	Wins           int    `json:"wins"`
	Losses         int    `json:"losses"`
}

// UserSegment struct handles the MongoDB schema for each users segments
type UserSegment struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// User struct handles the MongoDB schema for a user
type User struct {
	ID        int           `json:"_id"`
	FirstName string        `json:"firstName"`
	LastName  string        `json:"lastName"`
	FullName  string        `json:"fullName"`
	Token     string        `json:"token"`
	Photo     string        `json:"photo"`
	Email     string        `json:"email"`
	Friends   []Friend      `json:"friends"`
	Segments  []UserSegment `json:"segments"`
	Wins      int           `json:"wins"`
	Losses    int           `json:"losses"`
}

// GetUserByID gets a single stored user from MongoDB
func GetUserByID(id int) (*User, error) {
	var u User

	s, err := New()
	if err != nil {
		log.WithError(err).Error("Unable to create new MongoDB Session")
		return nil, err
	}

	defer s.Close()

	err = s.DB("heroku_zgxbr4j2").C("users").Find(bson.M{"_id": id}).One(&u)
	if err != nil {
		log.WithField("ID", id).Error("Unable to find user with id")
		return nil, err
	}

	return &u, err
}
