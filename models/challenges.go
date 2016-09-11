package models

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// UserChallengeEffort struct handles the MongoDB schema for each users challenge effort
type UserChallengeEffort struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Photo            string  `json:"photo"`
	Time             int     `json:"time"`
	Completed        bool    `json:"completed"`
	AverageCadence   float32 `json:"averageCadence"`
	AverageWatts     float32 `json:"averageWatts"`
	AverageHeartRate float32 `json:"averageHeartRate"`
	MaxHeartRate     int     `json:"maxHeartRate"`
}

// Challenge struct handles the MongoDB schema for a challenge
type Challenge struct {
	ID         bson.ObjectId       `bson:"_id,omitempty" json:"_id"`
	Segment    Segment             `json:"segment"`
	Challenger UserChallengeEffort `json:"challenger"`
	Challengee UserChallengeEffort `json:"challengee"`
	Status     string              `json:"status"`
	Created    time.Time           `json:"created"`
	Expires    time.Time           `json:"expires"`
	Completed  time.Time           `json:"completed"`
	Expired    bool                `json:"expired"`
	WinnerID   int                 `json:"winnerId"`
	WinnerName string              `json:"winnerName"`
	LoserID    int                 `json:"loserId"`
	LoserName  string              `json:"loserName"`
}

// GetChallengeByID gets a single stored challenge from MongoDB
func GetChallengeByID(id bson.ObjectId) (*Challenge, error) {
	var c Challenge

	s, err := New()
	if err != nil {
		log.WithError(err).Error("Unable to create new MongoDB Session")
		return nil, err
	}

	defer s.Close()

	err = s.DB("heroku_zgxbr4j2").C("challenges").Find(bson.M{"_id": id}).One(&c)
	if err != nil {
		log.WithField("ID", id).Error("Unable to find challenge with id")
		return nil, err
	}

	return &c, nil
}
