package handlers

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/pressly/chi"
	"gopkg.in/mgo.v2/bson"

	"github.com/jrzimmerman/bestrida-server-go/models"
)

// GetChallengeByID returns challenge by ID from the database
func GetChallengeByID(w http.ResponseWriter, r *http.Request) {
	res := New(w)
	defer res.Render()

	id := chi.URLParam(r, "id")
	// validate Challenge ID is a bson Object ID
	if !bson.IsObjectIdHex(id) {
		res.SetResponse(500, "Challenge ID cannot be converted to BSON Object ID")
		return
	}

	oid := bson.ObjectIdHex(id)
	// log challenge ID
	log.WithField("id", oid).Info("looking for challenge by ID")

	challenge, err := models.GetChallengeByID(oid)
	log.Infoln(challenge)

	if err != nil {
		log.WithField("ID", id).Debug("unable to get challenge by ID")
		res.SetResponse(500, err)
		return
	}

	res.SetResponse(200, challenge)
}
