package models

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// Effort struct handles the MongoDB schema for each users challenge effort
type Effort struct {
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
	ID         bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Segment    Segment       `bson:"segment" json:"segment"`
	Challenger Effort        `bson:"challenger" json:"challenger"`
	Challengee Effort        `bson:"challengee" json:"challengee"`
	Status     string        `bson:"status" json:"status"`
	Created    time.Time     `bson:"created" json:"created"`
	Expires    time.Time     `bson:"expires" json:"expires"`
	Completed  time.Time     `bson:"completed" json:"completed"`
	Expired    bool          `bson:"expired" json:"expired"`
	WinnerID   int           `bson:"winnerId" json:"winnerId"`
	WinnerName string        `bson:"winnerName" json:"winnerName"`
	LoserID    int           `bson:"loserId" json:"loserId"`
	LoserName  string        `bson:"loserName" json:"loserName"`
	CreatedAt  time.Time     `bson:"createdAt" json:"createdAt,omitempty"`
	UpdatedAt  time.Time     `bson:"updatedAt" json:"updatedAt,omitempty"`
	DeletedAt  *time.Time    `bson:"deletedAt" json:"deletedAt,omitempty"`
}

// GetChallengeByID gets a single stored challenge from MongoDB
func GetChallengeByID(id bson.ObjectId) (*Challenge, error) {
	var c Challenge

	if err := session.DB(name).C("challenges").Find(bson.M{"_id": id}).One(&c); err != nil {
		log.WithField("ID", id).Error("Unable to find challenge with id in database")
		return nil, err
	}

	return &c, nil
}

// CreateChallenge creates a new challenge in MongoDB
func CreateChallenge(c Challenge) error {
	if err := session.DB(name).C("challenges").Insert(c); err != nil {
		log.Errorf("Unable to create a new challenge:\n %v", err)
		return err
	}
	log.Printf("Challenge successfully created")
	return nil
}

// RemoveChallenge removes a challenge from MongoDB
func RemoveChallenge(id bson.ObjectId) error {
	if err := session.DB(name).C("challenges").RemoveId(id); err != nil {
		log.WithField("ID", id).Error("Unable to find challenge with id in database")
		return err
	}
	log.Printf("Challenge successfully removed: %v", id)
	return nil
}
