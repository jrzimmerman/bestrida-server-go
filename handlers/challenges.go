package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pressly/chi"
	"gopkg.in/mgo.v2/bson"

	"github.com/jrzimmerman/bestrida-server-go/models"
)

// GetChallengeByID returns challenge by ID from the database
func GetChallengeByID(w http.ResponseWriter, r *http.Request) {
	res := New(w)

	id := chi.URLParam(r, "id")
	// validate Challenge ID is a bson Object ID
	if !bson.IsObjectIdHex(id) {
		res.Render(http.StatusInternalServerError, "Challenge ID cannot be converted to BSON Object ID")
		return
	}

	oid := bson.ObjectIdHex(id)
	log.WithField("id", oid).Info("looking for challenge by ID")
	challenge, err := models.GetChallengeByID(oid)
	if err != nil {
		log.WithField("ID", id).Debug("unable to get challenge by ID")
		res.Render(http.StatusInternalServerError, err)
		return
	}

	res.Render(http.StatusOK, challenge)
}

type createRequest struct {
	SegmentID      *int    `json:"segmentId"`
	ChallengerID   *int    `json:"challengerId"`
	ChallengeeID   *int    `json:"challengeeId"`
	CompletionDate *string `json:"completionDate"`
}

// CreateChallenge creates a new challenge with post content
func CreateChallenge(w http.ResponseWriter, r *http.Request) {
	res := New(w)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		res.Render(http.StatusInternalServerError, "Could not read request body")
		return
	}

	var req createRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		res.Render(http.StatusInternalServerError, "Could not unmarshal request to create challenge")
		return
	}

	t := time.Now()
	created := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	e, err := time.Parse("Mon Jan 02 2006 15:04:05 GMT-0700 (MST)", *req.CompletionDate)
	if err != nil {
		res.Render(http.StatusInternalServerError, "Could not parse completion date from request")
		return
	}
	expires := time.Date(e.Year(), e.Month(), e.Day(), 23, 59, 59, 0, e.Location())

	challengeeUser, err := models.GetUserByID(int64(*req.ChallengeeID))
	if err != nil {
		log.WithField("ID", *req.ChallengeeID).Error("unable to retrieve challengee from database")
		res.Render(http.StatusInternalServerError, "unable to retrieve challengee from database")
		return
	}
	var challengee models.Opponent
	challengee.ID = challengeeUser.ID
	challengee.Name = challengeeUser.FullName
	challengee.Photo = challengeeUser.Photo
	challengee.Completed = false

	challengerUser, err := models.GetUserByID(int64(*req.ChallengerID))
	if err != nil {
		log.WithField("ID", *req.ChallengerID).Error("unable to retrieve challenger from database")
		res.Render(http.StatusInternalServerError, "unable to retrieve challenger from database")
		return
	}
	var challenger models.Opponent
	challenger.ID = challengerUser.ID
	challenger.Name = challengerUser.FullName
	challenger.Photo = challengerUser.Photo
	challenger.Completed = false

	segment, err := models.GetSegmentByID(int64(*req.SegmentID))
	if err != nil {
		log.WithField("id", *req.SegmentID).Debug("unable to get segment by ID")
		res.Render(http.StatusInternalServerError, err)
		return
	}

	var challenge models.Challenge
	challenge.ID = bson.NewObjectId()
	challenge.Challengee = &challengee
	challenge.Challenger = &challenger
	challenge.Segment = segment
	challenge.Status = "pending"
	challenge.Created = &created
	challenge.Expires = &expires
	challenge.CreatedAt = t
	challenge.UpdatedAt = t

	err = models.CreateChallenge(challenge)
	if err != nil {
		res.Render(http.StatusInternalServerError, "Could not create challenge in database")
		return
	}
	res.Render(http.StatusOK, challenge)
}

type updateRequest struct {
	ID bson.ObjectId `json:"id"`
}

// AcceptChallengeByID accepts a challenge
func AcceptChallengeByID(w http.ResponseWriter, r *http.Request) {
	res := New(w)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		res.Render(http.StatusInternalServerError, "Could not read request body")
		return
	}

	var req updateRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		res.Render(http.StatusInternalServerError, "Could not unmarshal request to update challenge")
		return
	}

	log.Infof("accepting challenge %v", req.ID)
	err = models.UpdateChallengeStatus(req.ID, "active", time.Now())
	if err != nil {
		res.Render(http.StatusInternalServerError, "Could not update challenge in database")
		return
	}
	res.Render(http.StatusOK, "challenge accepted")
}

// DeclineChallengeByID decline a challenge
func DeclineChallengeByID(w http.ResponseWriter, r *http.Request) {
	res := New(w)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		res.Render(http.StatusInternalServerError, "Could not read request body")
		return
	}

	var req updateRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		res.Render(http.StatusInternalServerError, "Could not unmarshal request to update challenge")
		return
	}

	log.Infof("declining challenge %v", req.ID)
	err = models.RemoveChallenge(req.ID)
	if err != nil {
		res.Render(http.StatusInternalServerError, "Could not remove challenge in database")
		return
	}
	res.Render(http.StatusOK, "challenge declined")
}

// GetPendingChallengesByUserID gets all pending challenges by user ID
func GetPendingChallengesByUserID(w http.ResponseWriter, r *http.Request) {
	res := New(w)

	id := chi.URLParam(r, "id")

	numID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.WithField("ID", numID).Error("unable to convert ID param")
		res.Render(http.StatusBadRequest, "unable to convert ID param")
		return
	}
	challenges, err := models.GetPendingChallenges(numID)
	if err != nil {
		log.WithField("ID", numID).Error("Could not retrieve pending challenges from database")
		res.Render(http.StatusInternalServerError, "Could not retrieve pending challenges from database")
		return
	}
	res.Render(http.StatusOK, challenges)
}

// GetActiveChallengesByUserID gets all active challenges by user ID
func GetActiveChallengesByUserID(w http.ResponseWriter, r *http.Request) {
	res := New(w)
	res.Render(http.StatusOK, "active challenges")
}

// GetCompletedChallengesByUserID gets all completed challenges by user ID
func GetCompletedChallengesByUserID(w http.ResponseWriter, r *http.Request) {
	res := New(w)
	res.Render(http.StatusOK, "completed challenges")
}
