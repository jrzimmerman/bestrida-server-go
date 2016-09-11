package handlers

import (
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"

	"github.com/jrzimmerman/bestrida-server-go/models"
)

// GetSegmentByID returns segment by ID from the database
func GetSegmentByID(c *gin.Context) {
	id := c.Param("id")
	numID, err := strconv.Atoi(id)
	if err != nil {
		log.WithField("id", numID).Debug("unable to convert ID param")
		c.JSON(500, err)
		return
	}
	segment, err := models.GetSegmentByID(numID)
	if err != nil {
		log.WithField("id", numID).Debug("unable to get segment by ID")
		c.JSON(500, err)
		return
	}

	c.JSON(200, segment)
}
