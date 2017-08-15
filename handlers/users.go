package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/jrzimmerman/bestrida-server-go/models"
	log "github.com/sirupsen/logrus"
)

// GetUserByID returns user by ID from the database
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	res := New(w)

	id := chi.URLParam(r, "id")

	numID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.WithField("ID", numID).Error("unable to convert ID param")
		res.Render(http.StatusInternalServerError, map[string]interface{}{"error": "unable to convert ID param"})
		return
	}

	user, err := models.GetUserByID(numID)
	if err != nil {
		log.WithField("ID", numID).Error("unable to get user by ID from database")
		res.Render(http.StatusInternalServerError, map[string]interface{}{"error": "unable to get user by ID from database"})
		return
	}

	log.WithField("USER ID", user.ID).Infof("user %d found", user.ID)
	res.Render(http.StatusOK, user)
}

// GetSegmentsByUserID returns a slice of user segments from the database
func GetSegmentsByUserID(w http.ResponseWriter, r *http.Request) {
	res := New(w)

	id := chi.URLParam(r, "id")

	numID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.WithField("ID", numID).Error("unable to convert ID param")
		res.Render(http.StatusInternalServerError, map[string]interface{}{"error": "unable to convert ID param"})
		return
	}

	user, err := models.GetUserByID(numID)
	if err != nil {
		log.WithField("ID", numID).Error("unable to get user by ID from database")
		res.Render(http.StatusInternalServerError, map[string]interface{}{"error": "unable to get user by ID from database"})
		return
	}

	log.WithField("USER ID", user.ID).Infof("user %d found", user.ID)
	log.WithField("USER ID", user.ID).Infof("user %d has %d segments", user.ID, len(user.Segments))
	res.Render(http.StatusOK, user.Segments)
}

// GetFriendsByUserID returns a slice of user segments from the database
func GetFriendsByUserID(w http.ResponseWriter, r *http.Request) {
	res := New(w)

	id := chi.URLParam(r, "id")

	numID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.WithField("ID", numID).Error("unable to convert ID param")
		res.Render(http.StatusInternalServerError, map[string]interface{}{"error": "unable to convert ID param"})
		return
	}

	user, err := models.GetUserByID(numID)
	if err != nil {
		log.WithField("USER ID", numID).Error("unable to get user by ID from database")
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "unable to get user by ID from database",
		})
		return
	}

	log.WithField("USER ID", user.ID).Infof("user %d found", user.ID)
	log.WithField("USER ID", user.ID).Infof("user %d has %d friends", user.ID, len(user.Friends))
	res.Render(http.StatusOK, user.Friends)
}
