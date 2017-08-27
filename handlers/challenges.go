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
	SegmentID      int        `json:"segmentId"`
	ChallengerID   int        `json:"challengerId"`
	ChallengeeID   int        `json:"challengeeId"`
	CompletionDate time.Time  `json:"completionDate"`
	CreationDate   *time.Time `json:"creationDate"`
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
	log.Infof("CreationDate: %v", req.CreationDate)
	var t time.Time
	if req.CreationDate != nil {
		t = *req.CreationDate
	} else {
		t = time.Now()
	}
	created := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	log.Infof("created date %v formatted successfully", created)

	e := req.CompletionDate
	expires := time.Date(e.Year(), e.Month(), e.Day(), 23, 59, 59, 0, e.Location())
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

// UpdateChallengeEffort grabs challenge effort information for a user from Strava
func UpdateChallengeEffort(ID bson.ObjectId, UserID int64) (*models.Challenge, error) {
	// Get challenge by ChallengeID from DB
	c, err := models.GetChallengeByID(ID)
	if err != nil {
		log.Errorf("unable to find challenge %v in DB", ID)
		return nil, err
	}
	// Get user by UserID from DB
	u, err := models.GetUserByID(UserID)
	if err != nil {
		log.Errorf("unable to find user %v in DB", UserID)
		return nil, err
	}

	// use our access token to grab generic segment info
	client := strava.NewClient(u.Token)

	// request efforts by segment ID between start and end dates
	log.Infof("Fetching segment %v info...", c.Segment.ID)
	log.Infof("beginning on %v", *c.Created)
	log.Infof("ending on %v", *c.Expires)
	efforts, err := strava.NewSegmentsService(client).ListEfforts(c.Segment.ID).AthleteId(u.ID).DateRange(*c.Created, *c.Expires).Do()
	if err != nil {

		return nil, err
	}
	log.Infof("rate limit percent: %v", strava.RateLimiting.FractionReached()*100)

	// check for segment efforts
	if len(efforts) > 0 {
		// update challenge in DB
		// efforts are returned sorted by time
		e := efforts[0]
		if e.Athlete.Id == c.Challengee.ID {
			// user is challengee
			c.Challengee.Time = &e.ElapsedTime
			c.Challengee.AverageCadence = &e.AverageCadence
			c.Challengee.AverageWatts = &e.AveragePower
			c.Challengee.AverageHeartRate = &e.AverageHeartrate
			c.Challengee.MaxHeartRate = &e.MaximumHeartrate
			c.Challengee.Completed = true
			c.UpdatedAt = time.Now()
			if err := c.UpdateChallenge(); err != nil {
				log.Error("unable to update challengee values in challenge")
				return nil, err
			}
			// check if challenger has completed to calculate winner
			if c.Challenger.Completed == true {
				// calculate winner
				log.Info("Opponent has completed as well, lets calculate a winner!")
			}
		} else if e.Athlete.Id == c.Challenger.ID {
			// user is challenger
			c.Challenger.Time = &e.ElapsedTime
			c.Challenger.AverageCadence = &e.AverageCadence
			c.Challenger.AverageWatts = &e.AveragePower
			c.Challenger.AverageHeartRate = &e.AverageHeartrate
			c.Challenger.MaxHeartRate = &e.MaximumHeartrate
			c.Challenger.Completed = true
			c.UpdatedAt = time.Now()
			if err := c.UpdateChallenge(); err != nil {
				log.Error("unable to update challenger values in challenge")
				return nil, err
			}
			// check if challengee has completed to calculate winner
			if c.Challengee.Completed == true {
				// calculate winner
				log.Info("Opponent has completed as well, lets calculate a winner!")
			}
		} else {
			log.Error("effort athlete id doesnt match challengee or challenger ID, something went wrong")
			return c, nil
		}
	} else {
		log.Info("No efforts returned from Strava")
		return nil, err
	}
	return c, nil
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

	c, err := UpdateChallengeEffort(req.ID, req.UserID)
	if c == nil || err != nil {
		log.Error("Could not update challenge effort")
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "Could not update challenge effort",
			"stack": err,
		})
		return
	}

	log.Infof("Challenge completed successfully")
	res.Render(http.StatusOK, c)
}

// UpdateChallengeResult determines the result of the challenge
func UpdateChallengeResult(id bson.ObjectId) error {
	c, err := models.GetChallengeByID(id)
	if err != nil {
		log.Errorf("challenge %v unable to be found in DB", id)
		return err
	}
	challengee, err := models.GetUserByID(c.Challengee.ID)
	if err != nil {
		log.Error("unable to find challengee")
		return err
	}
	challenger, err := models.GetUserByID(c.Challenger.ID)
	if err != nil {
		log.Error("unable to find challenger")
		return err
	}
	// completed refers to the time challenge was marked as complete
	completed := time.Now()
	if c.Challengee.Completed == false && c.Challenger.Completed == false {
		// remove challenge if no one completed
		models.RemoveChallenge(c.ID)
		if err != nil {
			log.Errorf("challenge %v unable to be removed from DB", c.ID)
			return err
		}
		return nil
	} else if c.Challengee.Completed == true && c.Challenger.Completed == false {
		// challengee was the only one who made an effort during the challenge
		c.WinnerID = &c.Challengee.ID
		c.WinnerName = &c.Challengee.Name
		c.LoserID = &c.Challenger.ID
		c.LoserName = &c.Challenger.Name
		c.Completed = &completed
		c.Status = "complete"
		c.Expired = true
		challengee.IncrementWins(challenger.ID)
		challenger.IncrementLosses(challengee.ID)
		challengee.IncrementSegments(c.Segment.ID)
		challenger.IncrementSegments(c.Segment.ID)
	} else if c.Challengee.Completed == false && c.Challenger.Completed == true {
		// challenger was the only one who made an effort during the challenge
		c.WinnerID = &c.Challenger.ID
		c.WinnerName = &c.Challenger.Name
		c.LoserID = &c.Challengee.ID
		c.LoserName = &c.Challengee.Name
		c.Completed = &completed
		c.Status = "complete"
		c.Expired = true
		challenger.IncrementWins(challengee.ID)
		challengee.IncrementLosses(challenger.ID)
		challenger.IncrementSegments(c.Segment.ID)
		challengee.IncrementSegments(c.Segment.ID)
	} else {
		// both challengers completed, determine winner based upon times
		if *c.Challengee.Time < *c.Challenger.Time {
			// challengee won
			c.WinnerID = &c.Challengee.ID
			c.WinnerName = &c.Challengee.Name
			c.LoserID = &c.Challenger.ID
			c.LoserName = &c.Challenger.Name
			c.Completed = &completed
			c.Status = "complete"
			c.Expired = true
			challengee.IncrementWins(challenger.ID)
			challenger.IncrementLosses(challengee.ID)
			challengee.IncrementSegments(c.Segment.ID)
			challenger.IncrementSegments(c.Segment.ID)
		} else if *c.Challenger.Time < *c.Challengee.Time {
			//  challenger won
			c.WinnerID = &c.Challenger.ID
			c.WinnerName = &c.Challenger.Name
			c.LoserID = &c.Challengee.ID
			c.LoserName = &c.Challengee.Name
			c.Completed = &completed
			c.Status = "complete"
			c.Expired = true
			challenger.IncrementWins(challengee.ID)
			challengee.IncrementLosses(challenger.ID)
			challenger.IncrementSegments(c.Segment.ID)
			challengee.IncrementSegments(c.Segment.ID)
		} else {
			// challenger and challengee times are the same
			log.Info("challenger and challengee effort times are the same")
			c.Completed = &completed
			c.Status = "complete"
			c.Expired = true
			challengee.IncrementSegments(c.Segment.ID)
			challenger.IncrementSegments(c.Segment.ID)
		}
	}
	if err := c.UpdateChallenge(); err != nil {
		log.Error("Unable to update challenge")
		return err
	}
	return nil
}

// CronComplete finds a list of expired challenges and processes them for completion
func CronComplete() {
	expired, err := models.GetExpiredChallenges()
	if err != nil {
		log.Error("Unable to find expired challenges")
		return
	}
	log.Infof("%d expired challenges returned from GetExpiredChallenges", len(*expired))
	for _, challenge := range *expired {
		log.Infof("challenge %v is expired on %v", challenge.ID, challenge.Expires)
		// only update challenge efforts if a challenge is not pending
		if challenge.Status != "pending" {
			// update efforts for both participants before determining a winner or loser
			UpdateChallengeEffort(challenge.ID, challenge.Challengee.ID)
			UpdateChallengeEffort(challenge.ID, challenge.Challenger.ID)
		}
		if err := UpdateChallengeResult(challenge.ID); err != nil {
			log.Error("Unable to update challenge result")
			return
		}
	}
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
