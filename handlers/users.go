package handlers

import (
	"strconv"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/jrzimmerman/bestrida-server-go/models"
)

// GetUserByID returns user by ID from the database
func GetUserByID(c *gin.Context) {
	id := c.Param("id")

	numID, err := strconv.Atoi(id)
	if err != nil {
		log.WithField("ID", numID).Error("unable to convert ID param")
		c.JSON(500, "unable to convert ID param")
		return
	}

	user, err := models.GetUserByID(numID)
	if err != nil {
		log.WithField("ID", numID).Error("unable to get user by ID")
		c.JSON(500, "unable to get user by ID")
		return
	}

	c.JSON(200, user)
}
