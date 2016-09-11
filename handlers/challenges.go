package handlers

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"

	"github.com/jrzimmerman/bestrida-server-go/models"
)

// GetChallengeByID returns challenge by ID from the database
func GetChallengeByID(c *gin.Context) {
	id := c.Param("id")
	oid := bson.ObjectIdHex(id)
	challenge, err := models.GetChallengeByID(oid)
	if err != nil {
		log.WithField("ID", id).Debug("unable to get challenge by ID")
		c.JSON(500, err)
		return
	}

	c.JSON(200, challenge)
}
