package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
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
	SegmentID      *int   `json:"segmentId"`
	ChallengerID   *int   `json:"challengerId"`
	ChallengeeID   *int   `json:"challengeeId"`
	CompletionDate *int64 `json:"completionDate"`
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
	log.Infof("SegmentID: %v", *req.SegmentID)
	log.Infof("ChallengerID: %v", *req.ChallengerID)
	log.Infof("ChallengeeID: %v", *req.ChallengeeID)
	log.Infof("CompletionDate: %v", *req.CompletionDate)

	t := time.Now()
	created := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	log.Infof("created date %v formatted successfully", created)

	e := time.Unix(*req.CompletionDate, 0)
	expires := time.Date(e.Year(), e.Month(), e.Day(), 23, 59, 59, 0, e.Location())
	log.Infof("expires date %v formatted successfully", expires)

	challengerUser, err := models.GetUserByID(int64(*req.ChallengerID))
	if err != nil {
		log.WithField("CHALLENGER ID", *req.ChallengerID).Error("unable to retrieve challenger from database")
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

	challengeeUser, err := models.GetUserByID(int64(*req.ChallengeeID))
	if err != nil {
		log.WithField("CHALLENGEE ID", *req.ChallengeeID).Error("unable to retrieve challengee from database")
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "unable to retrieve challengee from database",
			"stack": err,
		})
		return
	}
	challengee := models.Opponent{
		ID:        challengeeUser.ID,
		Name:      challengeeUser.FullName,
		Photo:     challengeeUser.Photo,
		Completed: false,
	}
	log.Infof("challengee %v formatted successfully", challengee.ID)

	segment, err := models.GetSegmentByID(int64(*req.SegmentID))
	if err != nil {
		log.WithField("SEGMENT ID", *req.SegmentID).Error("unable to get segment by ID")
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

// CompleteChallengeByID completes a challenge by challenge ID
func CompleteChallengeByID() {
	return
}

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
	challenges, err := models.GetAllChallenges(numID)
	if err != nil {
		log.WithField("ID", numID).Errorf("Could not retrieve challenges from database for user %v", numID)
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "Could not retrieve pending challenges from database",
			"stack": err,
		})
		return
	}
	res.Render(http.StatusOK, challenges)
}

// GetPendingChallengesByUserID gets all pending challenges by user ID
func GetPendingChallengesByUserID(w http.ResponseWriter, r *http.Request) {
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
	challenges, err := models.GetPendingChallenges(numID)
	if err != nil {
		log.WithField("ID", numID).Error("Could not retrieve pending challenges from database")
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "Could not retrieve pending challenges from database",
			"stack": err,
		})
		return
	}
	res.Render(http.StatusOK, challenges)
}

// GetActiveChallengesByUserID gets all active challenges by user ID
func GetActiveChallengesByUserID(w http.ResponseWriter, r *http.Request) {
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
	challenges, err := models.GetActiveChallenges(numID)
	if err != nil {
		log.WithField("ID", numID).Error("Could not retrieve active challenges from database")
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "Could not retrieve active challenges from database",
			"stack": err,
		})
		return
	}
	res.Render(http.StatusOK, challenges)
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
	challenges, err := models.GetCompletedChallenges(numID)
	if err != nil {
		log.WithField("ID", numID).Error("Could not retrieve active challenges from database")
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "Could not retrieve active challenges from database",
			"stack": err,
		})
		return
	}
	res.Render(http.StatusOK, challenges)
}
