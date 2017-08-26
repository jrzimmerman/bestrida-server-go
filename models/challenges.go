package models

import (
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// Opponent struct handles the database schema for each users challenge effort
type Opponent struct {
	ID               int64    `json:"id"`
	Name             string   `json:"name"`
	Photo            string   `json:"photo"`
	Completed        bool     `json:"completed"`
	Time             *int     `json:"time,omitempty"`
	AverageCadence   *float64 `json:"averageCadence,omitempty"`
	AverageWatts     *float64 `json:"averageWatts,omitempty"`
	AverageHeartRate *float64 `json:"averageHeartRate,omitempty"`
	MaxHeartRate     *float64 `json:"maxHeartRate,omitempty"`
}

// Challenge struct handles the database schema for a challenge
type Challenge struct {
	ID         bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Segment    *Segment      `bson:"segment" json:"segment"`
	Challenger *Opponent     `bson:"challenger" json:"challenger"`
	Challengee *Opponent     `bson:"challengee" json:"challengee"`
	Status     string        `bson:"status" json:"status"`
	Created    *time.Time    `bson:"created" json:"created,omitempty"`
	Expires    *time.Time    `bson:"expires" json:"expires,omitempty"`
	Completed  *time.Time    `bson:"completed" json:"completed,omitempty"`
	Expired    bool          `bson:"expired" json:"expired"`
	WinnerID   *int          `bson:"winnerId" json:"winnerId,omitempty"`
	WinnerName *string       `bson:"winnerName" json:"winnerName,omitempty"`
	LoserID    *int          `bson:"loserId" json:"loserId,omitempty"`
	LoserName  *string       `bson:"loserName" json:"loserName,omitempty"`
	CreatedAt  time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time     `bson:"updatedAt" json:"updatedAt"`
	DeletedAt  *time.Time    `bson:"deletedAt" json:"deletedAt,omitempty"`
}

// GetChallengeByID gets a single stored challenge from database
func GetChallengeByID(id bson.ObjectId) (*Challenge, error) {
	s := session.Copy()
	defer s.Close()

	var c Challenge

	if err := s.DB(name).C("challenges").Find(bson.M{"_id": id}).One(&c); err != nil {
		log.WithField("ID", id).Error("Unable to find challenge with id in database")
		return nil, err
	}

	return &c, nil
}

// CreateChallenge creates a new challenge in database
func CreateChallenge(c Challenge) error {
	s := session.Copy()
	defer s.Close()

	if err := s.DB(name).C("challenges").Insert(c); err != nil {
		log.WithField("CHALLENGE ID", c.ID).Errorf("Unable to create a new challenge:\n %v", err)
		return err
	}
	log.WithField("CHALLENGE ID", c.ID).Infof("Challenge %v successfully created", c.ID)
	return nil
}

// RemoveChallenge removes a challenge from database
func RemoveChallenge(id bson.ObjectId) error {
	s := session.Copy()
	defer s.Close()

	if err := s.DB(name).C("challenges").RemoveId(id); err != nil {
		log.WithField("CHALLENGE ID", id).Error("Unable to find challenge with id in database")
		return err
	}
	log.Infof("Challenge successfully removed: %v", id)
	return nil
}

// UpdateChallenge updates a challenge from database
func (c Challenge) UpdateChallenge() error {
	s := session.Copy()
	defer s.Close()

	if err := s.DB(name).C("challenges").UpdateId(c.ID, c); err != nil {
		log.WithField("CHALLENGE ID", c.ID).Errorf("Unable to update challenge %v in database", c.ID)
		return err
	}
	log.Infof("Challenge successfully updated: %v", c.ID)
	return nil
}

// UpdateChallengeStatus updates the challenge
func UpdateChallengeStatus(id bson.ObjectId, status string, updateTime time.Time) error {
	s := session.Copy()
	defer s.Close()

	if err := s.DB(name).C("challenges").Update(bson.M{"_id": id}, bson.M{"$set": bson.M{"status": status, "updatedAt": updateTime}}); err != nil {
		log.WithField("ID", id).Errorf("Unable to update challenge with id: %v in database", id)
		return err
	}
	log.Infof("Challenge successfully updated: %v", id)
	return nil
}

// GetAllChallenges get all challenges for a user from database
func GetAllChallenges(userID int64) (*[]Challenge, error) {
	s := session.Copy()
	defer s.Close()

	var challenges []Challenge
	err := s.DB(name).C("challenges").Find(bson.M{
		"$or": []bson.M{
			bson.M{"challengee.id": userID},
			bson.M{"challenger.id": userID},
		},
	}).Sort("expires").All(&challenges)
	if err != nil {
		log.WithField("ID", userID).Errorf("Unable to find challenges for user %d in database", userID)
		return nil, err
	}
	log.Infof("found %d challenges for user %v", len(challenges), userID)

	return &challenges, nil
}

// GetPendingChallenges get pending challenges by user ID from database
func GetPendingChallenges(userID int64) (*[]Challenge, error) {
	s := session.Copy()
	defer s.Close()

	var challenges []Challenge
	err := s.DB(name).C("challenges").Find(bson.M{
		"$or": []bson.M{
			bson.M{"challengee.id": userID, "status": "pending"},
			bson.M{"challenger.id": userID, "status": "pending"},
		},
	}).Sort("expires").All(&challenges)
	if err != nil {
		log.WithField("ID", userID).Errorf("Unable to find pending challenges for user %d in database", userID)
		return nil, err
	}
	log.Infof("found %d pending challenges for user %v", len(challenges), userID)

	return &challenges, nil
}

// GetActiveChallenges get active challenges by user ID from database
func GetActiveChallenges(userID int64) (*[]Challenge, error) {
	s := session.Copy()
	defer s.Close()

	var challenges []Challenge
	err := s.DB(name).C("challenges").Find(bson.M{
		"$or": []bson.M{
			bson.M{"challengee.id": userID, "challengee.completed": false, "status": "active"},
			bson.M{"challenger.id": userID, "challenger.completed": false, "status": "active"},
		},
	}).Sort("expires").All(&challenges)
	if err != nil {
		log.WithField("ID", userID).Errorf("Unable to find active challenges for user %d in database", userID)
		return nil, err
	}
	log.Infof("found %d active challenges for user %v", len(challenges), userID)

	return &challenges, nil
}

// GetCompletedChallenges get completed challenges by user ID from database
func GetCompletedChallenges(userID int64) (*[]Challenge, error) {
	s := session.Copy()
	defer s.Close()

	var challenges []Challenge
	err := s.DB(name).C("challenges").Find(bson.M{
		"$or": []bson.M{
			bson.M{"challengee.id": userID, "challengee.completed": true},
			bson.M{"challenger.id": userID, "challenger.completed": true},
		},
	}).Sort("updatedAt", "expires").All(&challenges)
	if err != nil {
		log.WithField("ID", userID).Errorf("Unable to find active challenges for user %d in database", userID)
		return nil, err
	}
	log.Infof("found %d completed challenges for user %v", len(challenges), userID)

	return &challenges, nil
}

// GetExpiredChallenges get expired challenges from database
func GetExpiredChallenges() (*[]Challenge, error) {
	s := session.Copy()
	defer s.Close()

	cutoff := time.Now()
	var challenges []Challenge
	err := s.DB(name).C("challenges").Find(bson.M{
		"expired": false,
		"expires": bson.M{"$lt": cutoff},
	}).All(&challenges)
	if err != nil {
		log.Errorf("Unable to find expired challenges in database")
		return nil, err
	}
	log.Infof("%d expired challenges found from DB", len(challenges))

	return &challenges, nil
}
