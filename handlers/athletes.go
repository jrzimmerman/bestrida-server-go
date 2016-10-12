package handlers

import (
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/jrzimmerman/bestrida-server-go/models"
	strava "github.com/strava/go.strava"
)

// GetAthleteByIDFromStrava returns the strava athlete with the specified ID
func GetAthleteByIDFromStrava(c *gin.Context) {
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

	client := strava.NewClient(user.Token)

	log.Info("Fetching athlete info...\n")
	athlete, err := strava.NewCurrentAthleteService(client).Get().Do()
	if err != nil {
		c.JSON(500, "Unable to retrieve athlete info")
		return
	}

	c.JSON(200, athlete)
}

// GetFriendsByUserIDFromStrava returns a list of friends for a specific user by ID from strava
func GetFriendsByUserIDFromStrava(c *gin.Context) {
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

	client := strava.NewClient(user.Token)

	log.Info("Fetching athlete friends info...\n")
	friends, err := strava.NewCurrentAthleteService(client).ListFriends().Do()
	if err != nil {
		c.JSON(500, "Unable to retrieve athlete friends")
		return
	}

	c.JSON(200, friends)
}

// GetSegmentsByUserIDFromStrava returns a list of segments for a specific user by ID from strava
func GetSegmentsByUserIDFromStrava(c *gin.Context) {
	id := c.Param("id")

	// convert id string to number
	numID, err := strconv.Atoi(id)
	if err != nil {
		log.WithField("ID", numID).Error("unable to convert ID param")
		c.JSON(500, "unable to convert ID param")
		return
	}

	// find user by numID to retrieve strava token
	user, err := models.GetUserByID(numID)
	if err != nil {
		log.WithField("ID", numID).Error("unable to retrieve user from database")
		c.JSON(500, "unable to retrieve user from database")
		return
	}

	// create new strava client with user token
	client := strava.NewClient(user.Token)

	log.Info("Fetching athlete activity summary info...\n")
	activities, err := strava.NewCurrentAthleteService(client).ListActivities().Do()
	if err != nil {
		c.JSON(500, "Unable to retrieve athlete activities summary")
		return
	}

	// range over activity summary to get activity details
	for _, activitySummary := range activities {
		log.WithField("activity summary", activitySummary).Info("activity summary")

		// request activity detail from strava to obtain segments
		activityDetail, err := strava.NewActivitiesService(client).Get(activitySummary.Id).Do()
		if err != nil {
			log.WithField("activity detail", activityDetail).Error("unable to retrieve activity detail")
			c.JSON(500, "Unable to retrieve activity detail")
			return
		}
		log.WithField("activity detail", activityDetail).Info("activity detail")

		// range over activity detail to obtain segment efforts
	}

	c.JSON(200, activities)
}
