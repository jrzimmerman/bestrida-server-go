package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/jrzimmerman/bestrida-server-go/models"
	log "github.com/sirupsen/logrus"
	"github.com/strava/go.strava"
)

// GetEffortsBySegmentIDFromStravaWithUserID returns efforts by segment ID from Strava
func GetEffortsBySegmentIDFromStravaWithUserID(w http.ResponseWriter, r *http.Request) {
	res := New(w)

	id := chi.URLParam(r, "id")

	numID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.WithField("ID", numID).Error("unable to convert ID param")
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "unable to convert segment ID param",
			"stack": err,
		})
		return
	}

	user, err := models.GetUserByID(numID)
	if err != nil {
		log.WithField("ID", numID).Error("unable to retrieve user from database")
		res.Render(http.StatusInternalServerError, map[string]interface{}{"error": "unable to retrieve user from database"})
		return
	}
	segmentID := chi.URLParam(r, "segmentID")
	numSegmentID, err := strconv.ParseInt(segmentID, 10, 64)
	if err != nil {
		log.WithField("Segment ID", numSegmentID).Debug("unable to convert Segment ID param")
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "unable to convert Segment ID param",
			"stack": err,
		})
		return
	}

	// use the users access token to grab segment effort info
	client := strava.NewClient(user.Token)

	log.Infof("Fetching segment %v info...", numSegmentID)
	efforts, err := strava.NewSegmentsService(client).ListEfforts(numSegmentID).AthleteId(user.ID).Do()
	if err != nil {
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "Unable to retrieve segment efforts info",
			"stack": err,
		})
		return
	}
	log.Infof("rate limit percent: %v", strava.RateLimiting.FractionReached()*100)

	res.Render(http.StatusOK, efforts)
}
