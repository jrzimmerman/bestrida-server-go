package models

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// Friend struct handles the MongoDB schema for each users friends
type Friend struct {
	ID             int    `bson:"id" json:"id"`
	Username       string `bson:"username" json:"username"`
	FirstName      string `bson:"firstname" json:"firstName"`
	LastName       string `bson:"lastname" json:"lastName"`
	FullName       string `bson:"fullName" json:"fullName"`
	Photo          string `bson:"photo" json:"photo"`
	ChallengeCount int    `bson:"challengeCount" json:"challengeCount"`
	Wins           int    `bson:"wins" json:"wins"`
	Losses         int    `bson:"losses" json:"losses"`
}

// UserSegment struct handles the MongoDB schema for each users segments
type UserSegment struct {
	ID    int    `bson:"_id" json:"id"`
	Name  string `bson:"name" json:"name"`
	Count int    `bson:"count" json:"count"`
}

// User struct handles the MongoDB schema for a user
type User struct {
	ID        int           `bson:"_id" json:"id"`
	FirstName string        `bson:"firstname" json:"firstName"`
	LastName  string        `bson:"lastname" json:"lastName"`
	FullName  string        `bson:"fullname" json:"fullName"`
	Token     string        `bson:"token" json:"token"`
	Photo     string        `bson:"photo" json:"photo"`
	Email     string        `bson:"email" json:"email"`
	Friends   []Friend      `bson:"friends" json:"friends"`
	Segments  []UserSegment `bson:"segments" json:"segments"`
	Wins      int           `bson:"wins" json:"wins"`
	Losses    int           `bson:"losses" json:"losses"`
}

// GetUserByID gets a single stored user from MongoDB
func GetUserByID(id int) (*User, error) {
	var u User

	if err := session.DB(name).C("users").Find(bson.M{"_id": id}).One(&u); err != nil {
		log.WithField("ID", id).Error("Unable to find user with id")
		return nil, err
	}

	return &u, nil
}

// ModifySegmentCount will modify a segment count by the count param for a specific user
func (u User) ModifySegmentCount(segmentID int, count int) error {

	for i := range u.Segments {
		if u.Segments[i].ID == segmentID {
			u.Segments[i].Count = u.Segments[i].Count + count
			break
		}
	}

	if err := session.DB(name).C("users").UpdateId(u.ID, &u); err != nil {
		return err
	}

	return nil
}
