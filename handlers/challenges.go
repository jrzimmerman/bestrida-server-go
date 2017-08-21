package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	strava "github.com/strava/go.strava"
	"gopkg.in/mgo.v2/bson"

	"github.com/jrzimmerman/bestrida-server-go/models"
)

// GetChallengeByID returns challenge by ID from the database
func GetChallengeByID(w http.ResponseWriter, r *http.Request) {
	res := New(w)

	id := chi.URLParam(r, "id")
	// validate Challenge ID is a bson Object ID
	if !bson.IsObjectIdHex(id) {
		res.Render(http.StatusInternalServerError, map[string]interface{}{"error": "Challenge ID cannot be converted to BSON Object ID"})
		return
	}

	oid := bson.ObjectIdHex(id)
	log.WithField("id", oid).Info("looking for challenge by ID")
	challenge, err := models.GetChallengeByID(oid)
	if err != nil {
		log.WithField("ID", id).Debug("unable to get challenge by ID")
		res.Render(http.StatusInternalServerError, map[string]interface{}{"error": err})
		return
	}

	res.Render(http.StatusOK, challenge)
}

type createRequest struct {
	SegmentID      int   `json:"segmentId"`
	ChallengerID   int   `json:"challengerId"`
	ChallengeeID   int   `json:"challengeeId"`
	CompletionDate int64 `json:"completionDate"`
}

// CreateChallenge creates a new challenge with post content
func CreateChallenge(w http.ResponseWriter, r *http.Request) {
	res := New(w)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "Could not read request body",
			"stack": err,
		})
		return
	}

	var req createRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "Could not unmarshal request to create challenge",
			"stack": err,
		})
		return
	}
	log.Infof("SegmentID: %v", req.SegmentID)
	log.Infof("ChallengerID: %v", req.ChallengerID)
	log.Infof("ChallengeeID: %v", req.ChallengeeID)
	log.Infof("CompletionDate: %v", req.CompletionDate)

	t := time.Now().UTC()
	created := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	log.Infof("created date %v formatted successfully", created)

	e := time.Unix(req.CompletionDate, 0).UTC()
	expires := time.Date(e.Year(), e.Month(), e.Day(), 23, 59, 59, 0, time.UTC)
	log.Infof("expires date %v formatted successfully", expires)

	challengerUser, err := models.GetUserByID(int64(req.ChallengerID))
	if err != nil {
		log.WithField("CHALLENGER ID", req.ChallengerID).Error("unable to retrieve challenger from database")
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "unable to retrieve challenger from database",
			"stack": err,
		})
		return
	}
	challenger := models.Opponent{
		ID:        challengerUser.ID,
		Name:      challengerUser.FullName,
		Photo:     challengerUser.Photo,
		Completed: false,
	}
	log.Infof("challenger %v formatted successfully", challenger.ID)

	var challengee models.Opponent
	for _, friend := range challengerUser.Friends {
		if friend.ID == int64(req.ChallengeeID) {
			log.Info("friend found")
			challengee = models.Opponent{
				ID:        friend.ID,
				Name:      friend.FullName,
				Photo:     friend.Photo,
				Completed: false,
			}
			break
		}
	}
	if challengee.ID == 0 {
		log.WithField("CHALLENGEE ID", req.ChallengeeID).Error("unable to retrieve challengee from database")
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "unable to retrieve challengee from database",
			"stack": err,
		})
		return
	}
	log.Infof("challengee %v formatted successfully", challengee.ID)

	segment, err := models.GetSegmentByID(int64(req.SegmentID))
	if err != nil {
		log.WithField("SEGMENT ID", req.SegmentID).Error("unable to get segment by ID")
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "unable to get segment by ID",
			"stack": err,
		})
		return
	}
	log.Infof("segment %v found from DB", segment.ID)

	challenge := models.Challenge{
		ID:         bson.NewObjectId(),
		Challengee: &challengee,
		Challenger: &challenger,
		Segment:    segment,
		Status:     "pending",
		Created:    &created,
		Expires:    &expires,
		CreatedAt:  t,
		UpdatedAt:  t,
	}
	log.WithField("CHALLENGE ID", challenge.ID).Infof("challenge %v formatted successfully", challenge.ID)

	err = models.CreateChallenge(challenge)
	if err != nil {
		log.Error("Could not create challenge in database")
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "Could not create challenge in database",
			"stack": err,
		})
		return
	}
	res.Render(http.StatusOK, challenge)
}

type completeRequest struct {
	ID     bson.ObjectId `json:"id"`
	UserID int64         `json:"userId"`
}

// CompleteChallengeByID completes a challenge by challenge ID
func CompleteChallengeByID(w http.ResponseWriter, r *http.Request) {
	res := New(w)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Could not ready request body")
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "Could not read request body",
			"stack": err,
		})
		return
	}

	var req completeRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Error("Could not unmarshal request to complete challenge")
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "Could not unmarshal request to complete challenge",
			"stack": err,
		})
		return
	}
	log.Infof("ChallengeID: %v", req.ID)
	log.Infof("UserID: %v", req.UserID)
	// Get challenge by ChallengeID
	challenge, err := models.GetChallengeByID(req.ID)
	if err != nil {
		log.Errorf("unable to find challenge %v in DB", req.ID)
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "Unable to find challenge in DB",
			"stack": err,
		})
		return
	}

	// use our access token to grab generic segment info
	client := strava.NewClient(accessToken)

	// request efforts by segment ID between start and end dates
	log.Infof("Fetching segment %v info...", challenge.Segment.ID)
	log.Infof("beginning on %v", *challenge.Created)
	log.Infof("ending on %v", *challenge.Expires)
	efforts, err := strava.NewSegmentsService(client).ListEfforts(challenge.Segment.ID).AthleteId(req.UserID).DateRange(*challenge.Created, *challenge.Expires).Do()
	if err != nil {
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "Unable to retrieve segment efforts info",
			"stack": err,
		})
		return
	}
	log.Infof("rate limit percent: %v", strava.RateLimiting.FractionReached()*100)

	// update challenge in DB

	res.Render(http.StatusOK, efforts)
}

// function checkForWinner (challengeId, callback) {
//   Challenge
//   .find({ _id: challengeId })
//   .then(function (challenges) {
//     var challenge = challenges[0];
//     var winner;
//     // If challenge is complete
//     if (challenge.challengerTime && challenge.challengeeTime) {
//       if (challenge.challengerTime === challenge.challengeeTime) {
//         callback(null, 'Challenge resulted in a tie');
//       } else {
//         winner = challenge.challengerTime < challenge.challengeeTime ? 'challenger' : 'challengee';
//       }
//       // Update wins or losses and challenge count for both users
//       if (winner === 'challenger') {
//         Users.incrementWins(challenge.challengerId, challenge.challengeeId);
//         Users.incrementLosses(challenge.challengeeId, challenge.challengerId);
//         updateChallengeWinnerAndLoser(challengeId, challenge.challengerId,
//           challenge.challengerName, challenge.challengeeId, challenge.challengeeName, callback);
//       } else if (winner === 'challengee') {
//         Users.incrementWins(challenge.challengeeId, challenge.challengerId);
//         Users.incrementLosses(challenge.challengerId, challenge.challengeeId);
//         updateChallengeWinnerAndLoser(challengeId, challenge.challengeeId, challenge.challengeeName, challenge.challengerId, challenge.challengerName, callback);
//       }
//       // Updates challenge status to 'Complete'
//       Challenge.update({ _id: challengeId }, { status: 'complete' }, function (err, raw) {
//         if (err) {
//           callback('Error updating challenge status to \'Complete\'');
//         }
//       });
//     } else {
//       callback(null, 'Effort has been updated, waiting for other user to complete');
//     }
//   })
//   .catch(function(error) {
//     callback('Error checking for winner: ' + error)
//   });
// }

// function updateChallengeWinnerAndLoser (challengeId, winnerId, winnerName, loserId, loserName, cb) {
//   var completeDate = new Date();
//   Challenge.update({ _id: challengeId},
//     {
//       winnerId: winnerId,
//       winnerName: winnerName,
//       loserId: loserId,
//       loserName: loserName,
//       expired: true,
//       completed: completeDate.toISOString()
//     },
//     function (err) {
//       if (err) {
//         cb('Error updating challenge winner/loser:' + util.stringify(err), null);
//       } else {
//         cb(null, 'Challenge updated with winner and loser');
//       }
//     });
// }

// module.exports.complete = function (challenge, effort, callback) {
//   if (!challenge) {
//     callback('Challenge not sent to complete challenge model');
//     return;
//   }

//   if (!effort) {
//     callback('Effort not sent to complete challenge model');
//     return;
//   }

//   // Can refactor code to pass the challenger/challengee role of user when
//   // API is called to save us this extra request to the database
//   Challenge.find({ _id: challenge.id })
//   .then(function () {
//     var userRole = challenge.challengerId === effort.athlete.id ? 'challenger' : 'challengee';
//     if (userRole === 'challenger') {
//       Challenge.update({ _id: challenge.id },
//         {
//           challengerTime: effort.elapsed_time,
//           challengerCompleted: true,
//           challengerAvgCadence: effort.average_cadence || 0,
//           challengerAvgWatts: effort.average_watts || 0,
//           challengerAvgHeartrate: effort.average_heartrate || 0,
//           challengerMaxHeartRate: effort.max_heartrate || 0,
//           segmentDistance: effort.segment.distance,
//           segmentAverageGrade: effort.segment.average_grade,
//           segmentMaxGrade: effort.segment.maximum_grade,
//           segmentElevationHigh: effort.segment.elevation_high,
//           segmentElevationLow: effort.segment.elevation_low,
//           segmentClimbCategory: effort.segment.climb_category
//         }, function (err, res) {
//         if (err) {
//           callback('Error updating challenge with user effort: ' + util.stringify(err));
//         } else {
//           console.log('Updated challenge with user effort: ' + !!res.nModified);
//         }
//       });
//     } else if (userRole === 'challengee') {
//       Challenge.update({ _id: challenge.id },
//         {
//           challengeeTime: effort.elapsed_time,
//           challengeeCompleted: true,
//           challengeeAvgCadence: effort.average_cadence || 0,
//           challengeeAvgWatts: effort.average_watts || 0,
//           challengeeAvgHeartrate: effort.average_heartrate || 0,
//           challengeeMaxHeartRate: effort.max_heartrate || 0,
//           segmentDistance: effort.segment.distance,
//           segmentAverageGrade: effort.segment.average_grade,
//           segmentMaxGrade: effort.segment.maximum_grade,
//           segmentElevationHigh: effort.segment.elevation_high,
//           segmentElevationLow: effort.segment.elevation_low,
//           segmentClimbCategory: effort.segment.climb_category
//         }, function (err, res) {
//         if (err) {
//           callback('Error updating challenge with user effort: ' + util.stringify(err));
//         } else {
//           console.log('Updated challenge with user effort: ' + !!res.nModified);
//         }
//       });
//     }
//   })
//   .then(function(){
//     // Checks if the challenge has a winner
//     checkForWinner(challenge.id, function(err, res) {
//       if (err) {
//         callback('Error checking for winner: ' + err);
//       } else {
//         callback(err, 'Successfully checked for winner: ' + res);
//       }
//     });
//   })
//   .catch(function(error) {
//     callback(error);
//   });
// };

// function updateChallengeResult(challenge) {
//   var completeDate = new Date();
//   Challenge.find({_id: challenge._id}, function(err, res){
//     if (res.length) {
//       var challenge = res[0];
//       // If challengee is the only user who completed challenge
//       if (challenge.challengeeCompleted) {
//         Challenge.update({ _id: challenge._id }, {
//           winnerId: challenge.challengeeId,
//           winnerName: challenge.challengeeName,
//           loserId: challenge.challengerId,
//           loserName: challenge.challengerName,
//           expired: true,
//           challengerCompleted: true,
//           status: 'complete',
//           completed: completeDate.toISOString()
//         }, function (err, raw) {
//           if (err) {
//             console.log(err);
//           }
//         });
//       // If challenger is the only user who completed challenge
//       } else if (challenge.challengerCompleted) {
//           Challenge.update({ _id: challenge._id }, {
//           winnerId: challenge.challengerId,
//           winnerName: challenge.challengerName,
//           loserId: challenge.challengeeId,
//           loserName: challenge.challengeeName,
//           expired: true,
//           challengeeCompleted: true,
//           status: 'complete',
//           completed: completeDate.toISOString()
//         }, function (err, raw) {
//           if (err) {
//             console.log(err);
//           }
//         });
//       }
//     }
//   });
// }

// module.exports.cronComplete = function() {
//   var cutoff = new Date();
//   var buffer = 0.5; // Buffer, in number of days
//   cutoff.setTime(cutoff.getTime() - buffer * 86400000);
//   Challenge.find({ expired: false, expires: { $lt: cutoff }})
//   .then(function(result){
//     console.log(result.length, 'expired challenges were found');

//     result.forEach(function(aChallenge) {
//       var onlyOneUserCompletedChallenge = (aChallenge.challengeeCompleted && !aChallenge.challengerCompleted) ||
//                                           (!aChallenge.challengeeCompleted && aChallenge.challengerCompleted);

//       // If neither completed challenge, delete challenge
//       if (!aChallenge.challengeeCompleted && !aChallenge.challengerCompleted) {
//         Challenge.find({ _id: aChallenge.id })
//         .remove(function(err, raw) {
//           if (err) console.log(err);
//           console.log('removed challenges: ', !!raw.nModified);
//         });

//       // Else if only one user completed challenge, set default winner
//       } else if (onlyOneUserCompletedChallenge) {
//         updateChallengeResult(aChallenge);
//       }
//     });
//   });
// };

type updateRequest struct {
	ID bson.ObjectId `json:"id"`
}

// AcceptChallengeByID accepts a challenge
func AcceptChallengeByID(w http.ResponseWriter, r *http.Request) {
	res := New(w)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "Could not read request body",
			"stack": err,
		})
		return
	}

	var req updateRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "Could not unmarshal request to update challenge",
			"stack": err,
		})
		return
	}

	log.Infof("accepting challenge %v", req.ID)
	err = models.UpdateChallengeStatus(req.ID, "active", time.Now())
	if err != nil {
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "Could not update challenge in database",
			"stack": err,
		})
		return
	}
	res.Render(http.StatusOK, "challenge accepted")
}

// DeclineChallengeByID decline a challenge
func DeclineChallengeByID(w http.ResponseWriter, r *http.Request) {
	res := New(w)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "Could not read request body",
			"stack": err,
		})
		return
	}

	var req updateRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "Could not unmarshal request to update challenge",
			"stack": err,
		})
		return
	}

	log.Infof("declining challenge %v", req.ID)
	err = models.RemoveChallenge(req.ID)
	if err != nil {
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "Could not remove challenge in database",
			"stack": err,
		})
		return
	}
	res.Render(http.StatusOK, "challenge declined")
}

// GetAllChallengesByUserID gets all pending challenges by user ID
func GetAllChallengesByUserID(w http.ResponseWriter, r *http.Request) {
	res := New(w)

	id := chi.URLParam(r, "id")

	numID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.WithField("ID", numID).Error("unable to convert user ID param")
		res.Render(http.StatusBadRequest, map[string]interface{}{
			"error": "unable to convert user ID param",
			"stack": err,
		})
		return
	}
	allChallenges, err := models.GetAllChallenges(numID)
	if err != nil {
		log.WithField("ID", numID).Errorf("Could not retrieve challenges from database for user %v", numID)
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "Could not retrieve pending challenges from database",
			"stack": err,
		})
		return
	}
	res.Render(http.StatusOK, allChallenges)
}

// GetPendingChallengesByUserID gets all pending challenges by user ID
func GetPendingChallengesByUserID(w http.ResponseWriter, r *http.Request) {
	res := New(w)

	id := chi.URLParam(r, "id")

	numID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.WithField("USER ID", numID).Error("unable to convert user ID param")
		res.Render(http.StatusBadRequest, map[string]interface{}{
			"error": "unable to convert user ID param",
			"stack": err,
		})
		return
	}
	pendingChallenges, err := models.GetPendingChallenges(numID)
	if err != nil {
		log.WithField("USER ID", numID).Error("Could not retrieve pending challenges from database")
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "Could not retrieve pending challenges from database",
			"stack": err,
		})
		return
	}
	res.Render(http.StatusOK, pendingChallenges)
}

// GetActiveChallengesByUserID gets all active challenges by user ID
func GetActiveChallengesByUserID(w http.ResponseWriter, r *http.Request) {
	res := New(w)

	id := chi.URLParam(r, "id")

	numID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.WithField("USER ID", numID).Error("unable to convert user ID param")
		res.Render(http.StatusBadRequest, map[string]interface{}{
			"error": "unable to convert user ID param",
			"stack": err,
		})
		return
	}
	activeChallenges, err := models.GetActiveChallenges(numID)
	if err != nil {
		log.WithField("ID", numID).Error("Could not retrieve active challenges from database")
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "Could not retrieve active challenges from database",
			"stack": err,
		})
		return
	}
	res.Render(http.StatusOK, activeChallenges)
}

// GetCompletedChallengesByUserID gets all completed challenges by user ID
func GetCompletedChallengesByUserID(w http.ResponseWriter, r *http.Request) {
	res := New(w)

	id := chi.URLParam(r, "id")

	numID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.WithField("ID", numID).Error("unable to convert user ID param")
		res.Render(http.StatusBadRequest, map[string]interface{}{
			"error": "unable to convert user ID param",
			"stack": err,
		})
		return
	}
	completedChallenges, err := models.GetCompletedChallenges(numID)
	if err != nil {
		log.WithField("ID", numID).Error("Could not retrieve active challenges from database")
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "Could not retrieve active challenges from database",
			"stack": err,
		})
		return
	}
	res.Render(http.StatusOK, completedChallenges)
}
