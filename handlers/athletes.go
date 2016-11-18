package handlers

import (
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/jrzimmerman/bestrida-server-go/models"
	strava "github.com/strava/go.strava"
	"gopkg.in/gin-gonic/gin.v1"
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
	// retrieve a list of users segments from Strava API
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
	// retrieve a list of users friends from Strava API
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
	var segments []*strava.SegmentDetailed

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
	// the activity summary does not have segment efforts to view recent segments
	for _, activitySummary := range activities {
		log.WithFields(map[string]interface{}{
			"NAME": activitySummary.Name,
			"ID":   activitySummary.Id,
		}).Info("activity summary")

		// request activity detail from strava to obtain segments
		activityDetail, err := strava.NewActivitiesService(client).Get(activitySummary.Id).Do()
		if err != nil {
			log.WithFields(map[string]interface{}{
				"NAME": activityDetail.Name,
				"ID":   activityDetail.Id,
			}).Error("unable to retrieve activity detail")
			c.JSON(500, map[string]interface{}{
				"error":    "Unable to retrieve activity detail",
				"activity": activityDetail,
			})
			return
		}

		// range over segment efforts from the activity detail
		// to obtain segment details to cache
		for _, effort := range activityDetail.SegmentEfforts {
			log.WithField("SEGMENT", effort.Name).Info("segment effort from activity detail")
			segmentDetail, err := strava.NewSegmentsService(client).Get(effort.Segment.Id).Do()
			if err != nil {
				log.WithFields(map[string]interface{}{
					"NAME": effort.Name,
					"ID":   effort.Id,
				}).Error("unable to retrieve activity detail")
				c.JSON(500, map[string]interface{}{
					"error":   "Unable to retrieve activity detail",
					"segment": effort,
				})
				return
			}
			log.WithField("SEGMENT DETAIL", segmentDetail).Info("retrieved segment detail from strava")
			segments = append(segments, segmentDetail)
		}
	}
	c.JSON(200, segments)
}
