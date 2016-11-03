package handlers

import (
	"strconv"

	log "github.com/Sirupsen/logrus"
	strava "github.com/strava/go.strava"
	"gopkg.in/gin-gonic/gin.v1"

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

	// log segment ID
	log.WithField("id", numID).Info("looking for segment by ID")

	segment, err := models.GetSegmentByID(numID)
	if err != nil {
		log.WithField("id", numID).Debug("unable to get segment by ID")
		c.JSON(500, err)
		return
	}

	c.JSON(200, segment)
}

// GetSegmentByIDFromStrava returns the strava segment with the specified ID
func GetSegmentByIDFromStrava(c *gin.Context) {
	id := c.Param("id")

	numID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.WithField("ID", numID).Error("unable to convert segment ID param")
		c.JSON(500, "unable to convert segment ID param")
		return
	}

	// use our access token to grab generic segment info
	client := strava.NewClient(accessToken)

	log.Infof("Fetching segment %v info...\n", id)
	segment, err := strava.NewSegmentsService(client).Get(numID).Do()
	if err != nil {
		c.JSON(500, "Unable to retrieve segment info")
		return
	}

	c.JSON(200, segment)
}
