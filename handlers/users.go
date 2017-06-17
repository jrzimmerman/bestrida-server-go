package handlers

import (
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/pressly/chi"

	"github.com/jrzimmerman/bestrida-server-go/models"
)

// GetUserByID returns user by ID from the database
func GetUserByID(w http.ResponseWriter, r *http.Request) {
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
		log.WithField("ID", numID).Error("unable to get user by ID")
		res.Render(500, "unable to get user by ID")
		return
	}

	log.WithField("USER ID", user.ID).Infof("user %d found", user.ID)
	res.Render(200, user)
}
