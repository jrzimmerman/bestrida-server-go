package handlers

import (
	"strconv"

	log "github.com/Sirupsen/logrus"
	strava "github.com/strava/go.strava"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/jrzimmerman/bestrida-server-go/models"
)

// GetEffortsBySegmentIDFromStrava returns efforts by segment ID from Strava
func GetEffortsBySegmentIDFromStrava(c *gin.Context) {
	id := c.Param("id")
	numID, err := strconv.Atoi(id)
	if err != nil {
		log.WithField("ID", numID).Error("unable to convert ID param")
		c.JSON(500, "unable to convert ID param")
		return
	}

	user, err := models.GetUserByID(numID)
	if err != nil {
		log.WithField("ID", numID).Error("unable to retrieve user from database")
		c.JSON(500, "unable to retrieve user from database")
		return
	}
	userID := int64(user.ID)

	segmentID := c.Param("segmentID")
	numSegmentID, err := strconv.ParseInt(segmentID, 10, 64)
	if err != nil {
		log.WithField("Segment ID", numSegmentID).Debug("unable to convert Segment ID param")
		c.JSON(500, err)
		return
	}

	// use our access token to grab generic segment info
	client := strava.NewClient(user.Token)

	log.Infof("Fetching segment %v info...", numSegmentID)
	efforts, err := strava.NewSegmentsService(client).ListEfforts(numSegmentID).AthleteId(userID).Do()
	if err != nil {
		c.JSON(500, "Unable to retrieve segment efforts info")
		return
	}

	c.JSON(200, efforts)
}
