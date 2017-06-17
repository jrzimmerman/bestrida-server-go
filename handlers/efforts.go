package handlers

import (
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/jrzimmerman/bestrida-server-go/models"
	"github.com/pressly/chi"
	strava "github.com/strava/go.strava"
)

// GetEffortsBySegmentIDFromStrava returns efforts by segment ID from Strava
func GetEffortsBySegmentIDFromStrava(w http.ResponseWriter, r *http.Request) {
	res := New(w)

	id := chi.URLParam(r, "id")

	numID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.WithField("ID", numID).Error("unable to convert ID param")
		res.Render(500, "unable to convert ID param")
		return
	}

	user, err := models.GetUserByID(numID)
	if err != nil {
		log.WithField("ID", numID).Error("unable to retrieve user from database")
		res.Render(500, "unable to retrieve user from database")
		return
	}
	userID := int64(user.ID)

	segmentID := chi.URLParam(r, "segmentID")
	numSegmentID, err := strconv.ParseInt(segmentID, 10, 64)
	if err != nil {
		log.WithField("Segment ID", numSegmentID).Debug("unable to convert Segment ID param")
		res.Render(500, err)
		return
	}

	// use our access token to grab generic segment info
	client := strava.NewClient(user.Token)

	log.Infof("Fetching segment %v info...", numSegmentID)
	efforts, err := strava.NewSegmentsService(client).ListEfforts(numSegmentID).AthleteId(userID).Do()
	if err != nil {
		res.Render(500, "Unable to retrieve segment efforts info")
		return
	}

	res.Render(200, efforts)
}
