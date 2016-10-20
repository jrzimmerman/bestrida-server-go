package handlers

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2/bson"

	"github.com/jrzimmerman/bestrida-server-go/models"
)

// GetChallengeByID returns challenge by ID from the database
func GetChallengeByID(c *gin.Context) {
	id := c.Param("id")
	// validate Challenge ID is a bson Object ID
	if !bson.IsObjectIdHex(id) {
		c.JSON(500, "Challenge ID cannot be converted to BSON Object ID")
		return
	}

	oid := bson.ObjectIdHex(id)
	// log challenge ID
	log.WithField("id", oid).Info("looking for challenge by ID")

	challenge, err := models.GetChallengeByID(oid)
	log.Infoln(challenge)

	if err != nil {
		log.WithField("ID", id).Debug("unable to get challenge by ID")
		c.JSON(500, err)
		return
	}

	c.JSON(200, challenge)
}
