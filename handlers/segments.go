package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/jrzimmerman/bestrida-server-go/models"
	log "github.com/sirupsen/logrus"
	"github.com/strava/go.strava"
)

// GetSegmentByID returns segment by ID from the database
func GetSegmentByID(w http.ResponseWriter, r *http.Request) {
	res := New(w)

	id := chi.URLParam(r, "id")
	numID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.WithField("id", numID).Debug("unable to convert segment ID param")
		res.Render(http.StatusInternalServerError, map[string]interface{}{"error": "unable to convert segment ID param"})
		return
	}

	log.WithField("id", numID).Info("looking for segment by ID")
	segment, err := models.GetSegmentByID(numID)
	if err != nil {
		log.WithField("id", numID).Debug("unable to retrieve segment by ID from database")
		res.Render(http.StatusInternalServerError, map[string]interface{}{"error": "unable to retrieve segment by ID from database"})
		return
	}

	res.Render(http.StatusOK, segment)
}

// GetSegmentByIDFromStrava returns the strava segment with the specified ID
func GetSegmentByIDFromStrava(w http.ResponseWriter, r *http.Request) {
	res := New(w)

	id := chi.URLParam(r, "id")

	numID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.WithField("ID", numID).Error("unable to convert segment ID param")
		res.Render(http.StatusInternalServerError, map[string]interface{}{"error": "unable to convert segment ID param"})
		return
	}

	// use our access token to grab generic segment info
	client := strava.NewClient(accessToken)

	log.Infof("Fetching segment %v info from strava...", id)
	segment, err := strava.NewSegmentsService(client).Get(numID).Do()
	if err != nil {
		res.Render(http.StatusInternalServerError, map[string]interface{}{"error": "Unable to retrieve segment info"})
		return
	}
	log.Infof("rate limit percent: %v", strava.RateLimiting.FractionReached()*100)
	log.Infof("segment %v retrieved from strava", segment.Id)
	res.Render(http.StatusOK, segment)

	// store segment in background after render
	s, err := models.GetSegmentByID(segment.Id)
	if err != nil {
		models.SaveSegment(segment)
	} else {
		s.UpdateSegment(segment)
	}
}
