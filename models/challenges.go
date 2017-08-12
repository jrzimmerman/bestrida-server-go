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
	AverageCadence   *float32 `json:"averageCadence,omitempty"`
	AverageWatts     *float32 `json:"averageWatts,omitempty"`
	AverageHeartRate *float32 `json:"averageHeartRate,omitempty"`
	MaxHeartRate     *int     `json:"maxHeartRate,omitempty"`
}

// Challenge struct handles the database schema for a challenge
type Challenge struct {
	ID         bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Segment    *Segment      `bson:"segment" json:"segment,omitempty"`
	Challenger *Opponent     `bson:"challenger" json:"challenger,omitempty"`
	Challengee *Opponent     `bson:"challengee" json:"challengee,omitempty"`
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
	var c Challenge

	if err := session.DB(name).C("challenges").Find(bson.M{"_id": id}).One(&c); err != nil {
		log.WithField("ID", id).Error("Unable to find challenge with id in database")
		return nil, err
	}

	return &c, nil
}

// CreateChallenge creates a new challenge in database
func CreateChallenge(c Challenge) error {
	if err := session.DB(name).C("challenges").Insert(c); err != nil {
		log.Errorf("Unable to create a new challenge:\n %v", err)
		return err
	}
	log.Infof("Challenge successfully created")
	return nil
}

// RemoveChallenge removes a challenge from database
func RemoveChallenge(id bson.ObjectId) error {
	if err := session.DB(name).C("challenges").RemoveId(id); err != nil {
		log.WithField("ID", id).Error("Unable to find challenge with id in database")
		return err
	}
	log.Infof("Challenge successfully removed: %v", id)
	return nil
}

// UpdateChallengeStatus updates the challenge
func UpdateChallengeStatus(id bson.ObjectId, status string, updateTime time.Time) error {
	if err := session.DB(name).C("challenges").Update(bson.M{"_id": id}, bson.M{"$set": bson.M{"status": status, "updatedAt": updateTime}}); err != nil {
		log.WithField("ID", id).Errorf("Unable to update challenge with id: %v in database", id)
		return err
	}
	log.Infof("Challenge successfully updated: %v", id)
	return nil
}

// GetPendingChallenges get pending challenges by user ID from database
func GetPendingChallenges(userID int64) (*[]Challenge, error) {
	var challenges []Challenge
	err := session.DB(name).C("challenges").Find(bson.M{
		"$or": []bson.M{
			bson.M{"challengee.id": userID, "status": "pending"},
			bson.M{"challenger.id": userID, "status": "pending"},
		},
	}).Sort("expires").All(&challenges)
	if err != nil {
		log.WithField("ID", userID).Errorf("Unable to find pending challenges for user %d in database", userID)
		return nil, err
	}
	log.Infof("found %d pending challenges", len(challenges))

	return &challenges, nil
}

// GetActiveChallenges get active challenges by user ID from database
func GetActiveChallenges(userID int64) (*[]Challenge, error) {
	var challenges []Challenge
	err := session.DB(name).C("challenges").Find(bson.M{
		"$or": []bson.M{
			bson.M{"challengee.id": userID, "challengee.completed": false, "status": "active"},
			bson.M{"challenger.id": userID, "challenger.completed": false, "status": "active"},
		},
	}).Sort("expires").All(&challenges)
	if err != nil {
		log.WithField("ID", userID).Errorf("Unable to find active challenges for user %d in database", userID)
		return nil, err
	}
	log.Infof("found %d active challenges", len(challenges))

	return &challenges, nil
}

// GetCompletedChallenges get completed challenges by user ID from database
func GetCompletedChallenges(userID int64) (*[]Challenge, error) {
	var challenges []Challenge
	err := session.DB(name).C("challenges").Find(bson.M{
		"$or": []bson.M{
			bson.M{"challengee.id": userID, "challengee.completed": true, "status": "complete"},
			bson.M{"challenger.id": userID, "challenger.completed": true, "status": "complete"},
		},
	}).Sort("updatedAt", "expires").All(&challenges)
	if err != nil {
		log.WithField("ID", userID).Errorf("Unable to find active challenges for user %d in database", userID)
		return nil, err
	}
	log.Infof("found %d completed challenges", len(challenges))

	return &challenges, nil
}

// module.exports.getChallenges = function (user, status, callback) {
//   if (!user) callback('No user defined');
//   if (!status) callback('No status defined');
//   if (status === 'complete') {
//     Challenge
//     .find({
//       $or: [
//         { challengerId: user, challengerCompleted: true },
//         { challengeeId: user, challengeeCompleted: true }
//       ],
//     })
//     .sort([['updatedAt','descending'],['expires','descending']])
//     .exec(function (err, challenges) {
//       if (err) {
//         callback('error finding completed challenges: ' + err);
//       } else {
//         callback(err, challenges);
//       }
//     });
//   } else if (status === 'active') {
//     Challenge
//     .find({
//       $or: [
//         { challengerId: user, challengerCompleted: false, status: status },
//         { challengeeId: user, challengeeCompleted: false, status: status }
//       ]
//     })
//     .sort({ expires: 'ascending' })
//     .exec(function (err, challenges) {
//       if (err) {
//         callback('error finding active challenges: ' + err);
//       } else {
//         callback(err, challenges);
//       }
//     });
//   } else if (status === 'pending') {
//     Challenge.find({
//       $or: [
//         { challengeeId: user, status: 'pending' },
//         { challengerId: user, status: 'pending' }
//       ]
//     })
//     .sort({ expires: 'ascending' })
//     .exec(function (err, challenges) {
//       if (err) {
//         callback('error finding pending challenges: ' + err);
//       } else {
//         callback(err, challenges);
//       }
//     });
//   } else {
//     callback('Error getting challenges');
//   }
// };
