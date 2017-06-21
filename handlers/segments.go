package handlers

import (
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/pressly/chi"
	strava "github.com/strava/go.strava"

	"github.com/jrzimmerman/bestrida-server-go/models"
)

// GetSegmentByID returns segment by ID from the database
func GetSegmentByID(w http.ResponseWriter, r *http.Request) {
	res := New(w)

	id := chi.URLParam(r, "id")
	numID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logrus.WithField("id", numID).Debug("unable to convert ID param")
		res.Render(500, err)
		return
	}

	// logrus.segment ID
	logrus.WithField("id", numID).Info("looking for segment by ID")

	segment, err := models.GetSegmentByID(numID)
	if err != nil {
		logrus.WithField("id", numID).Debug("unable to get segment by ID")
		res.Render(500, err)
		return
	}

	res.Render(200, segment)
}

// GetSegmentByIDFromStrava returns the strava segment with the specified ID
func GetSegmentByIDFromStrava(w http.ResponseWriter, r *http.Request) {
	res := New(w)

	id := chi.URLParam(r, "id")

	numID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logrus.WithField("ID", numID).Error("unable to convert segment ID param")
		res.Render(500, "unable to convert segment ID param")
		return
	}

	// use our access token to grab generic segment info
	client := strava.NewClient(accessToken)

	logrus.Infof("Fetching segment %v info from strava...", id)
	segment, err := strava.NewSegmentsService(client).Get(numID).Do()
	if err != nil {
		res.Render(500, "Unable to retrieve segment info")
		return
	}
	logrus.Infof("rate limit percent: %v", strava.RateLimiting.FractionReached()*100)
	logrus.Infof("segment %v retrieved from strava", segment.Id)
	res.Render(200, segment)

	// store segment in background after render
	s, err := models.GetSegmentByID(segment.Id)
	if err != nil {
		models.SaveSegment(segment)
	} else {
		s.UpdateSegment(segment)
	}
}
