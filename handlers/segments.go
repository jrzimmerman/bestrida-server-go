package handlers

import (
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/pressly/chi"
	strava "github.com/strava/go.strava"

	"github.com/jrzimmerman/bestrida-server-go/models"
)

// GetSegmentByID returns segment by ID from the database
func GetSegmentByID(w http.ResponseWriter, r *http.Request) {
	res := New(w)
	defer res.Render()

	id := chi.URLParam(r, "id")
	numID, err := strconv.Atoi(id)
	if err != nil {
		log.WithField("id", numID).Debug("unable to convert ID param")
		res.SetResponse(500, err)
		return
	}

	// log segment ID
	log.WithField("id", numID).Info("looking for segment by ID")

	segment, err := models.GetSegmentByID(numID)
	if err != nil {
		log.WithField("id", numID).Debug("unable to get segment by ID")
		res.SetResponse(500, err)
		return
	}

	res.SetResponse(200, segment)
}

// GetSegmentByIDFromStrava returns the strava segment with the specified ID
func GetSegmentByIDFromStrava(w http.ResponseWriter, r *http.Request) {
	res := New(w)
	defer res.Render()

	id := chi.URLParam(r, "id")

	numID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.WithField("ID", numID).Error("unable to convert segment ID param")
		res.SetResponse(500, "unable to convert segment ID param")
		return
	}

	// use our access token to grab generic segment info
	client := strava.NewClient(accessToken)

	log.Infof("Fetching segment %v info...", id)
	segment, err := strava.NewSegmentsService(client).Get(numID).Do()
	if err != nil {
		res.SetResponse(500, "Unable to retrieve segment info")
		return
	}

	res.SetResponse(200, segment)
}
