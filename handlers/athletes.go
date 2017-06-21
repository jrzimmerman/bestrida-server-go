package handlers

import (
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/jrzimmerman/bestrida-server-go/models"
	"github.com/pressly/chi"
	strava "github.com/strava/go.strava"
)

// GetAthleteByIDFromStrava returns the strava athlete with the specified ID
func GetAthleteByIDFromStrava(w http.ResponseWriter, r *http.Request) {
	res := New(w)

	id := chi.URLParam(r, "id")

	numID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logrus.WithField("ID", numID).Error("unable to convert ID param")
		res.Render(400, "unable to convert ID param")
		return
	}

	user, err := models.GetUserByID(numID)
	if err != nil {
		logrus.WithField("ID", numID).Error("unable to retrieve user from database")
		res.Render(500, "unable to retrieve user from database")
		return
	}

	client := strava.NewClient(user.Token)

	logrus.Info("Fetching athlete info...\n")
	// retrieve a list of users segments from Strava API
	athlete, err := strava.NewCurrentAthleteService(client).Get().Do()
	if err != nil {
		res.Render(500, "Unable to retrieve athlete info")
		return
	}

	logrus.Infof("rate limit percent: %v", strava.RateLimiting.FractionReached()*100)
	logrus.Infof("athlete %v retrieved from strava", athlete.Id)
	u, err := user.UpdateAthlete(athlete)
	if err != nil {
		logrus.WithError(err).Errorf("unable to update athlete %d", athlete.Id)
	}
	res.Render(200, &u)
}

// GetFriendsByUserIDFromStrava returns a list of friends for a specific user by ID from strava
func GetFriendsByUserIDFromStrava(w http.ResponseWriter, r *http.Request) {
	res := New(w)

	id := chi.URLParam(r, "id")

	numID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logrus.WithField("ID", numID).Error("unable to convert ID param")
		res.Render(500, "unable to convert ID param")
		return
	}

	user, err := models.GetUserByID(numID)
	if err != nil {
		logrus.WithField("ID", numID).Error("unable to retrieve user from database")
		res.Render(500, "unable to retrieve user from database")
		return
	}

	client := strava.NewClient(user.Token)

	logrus.Info("Fetching athlete friends info...\n")
	// retrieve a list of users friends from Strava API
	friends, err := strava.NewCurrentAthleteService(client).ListFriends().Do()
	if err != nil {
		res.Render(500, "Unable to retrieve athlete friends")
		return
	}

	// store friends for a specific user

	// return request
	logrus.Infof("rate limit percent: %v", strava.RateLimiting.FractionReached()*100)
	logrus.Infof("found %d friends for athlete %d from strava", len(friends), user.ID)
	res.Render(200, friends)
}

// GetSegmentsByUserIDFromStrava returns a list of segments for a specific user by ID from strava
func GetSegmentsByUserIDFromStrava(w http.ResponseWriter, r *http.Request) {
	res := New(w)

	id := chi.URLParam(r, "id")

	var segments []*models.Segment

	// convert user id string from url param to number
	numID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logrus.WithField("USER ID", numID).Error("unable to convert USER ID param")
		res.Render(500, "unable to convert USER ID param")
		return
	}

	// find user by numID to retrieve strava token
	user, err := models.GetUserByID(numID)
	if err != nil {
		logrus.WithField("ID", numID).Error("unable to retrieve user from database")
		res.Render(404, "unable to retrieve user from database")
		return
	}

	// create new strava client with user token
	client := strava.NewClient(user.Token)

	logrus.Info("Fetching athlete activity summary info...\n")
	activities, err := strava.NewCurrentAthleteService(client).ListActivities().Page(1).PerPage(200).Do()
	if err != nil {
		res.Render(500, "Unable to retrieve athlete activities summary")
		return
	}

	// create a segment map to return only unique segments
	segmentMap := make(map[int64]*models.Segment, 0)

	// range over activity summary to get activity details
	// the activity summary does not have segment efforts to view recent segments
	for _, activitySummary := range activities {
		logrus.WithFields(map[string]interface{}{
			"NAME": activitySummary.Name,
			"ID":   activitySummary.Id,
		}).Info("activity summary")

		// request activity detail from strava to obtain segments
		activityDetail, err := strava.NewActivitiesService(client).Get(activitySummary.Id).Do()
		if err != nil {
			logrus.WithFields(map[string]interface{}{
				"NAME": activityDetail.Name,
				"ID":   activityDetail.Id,
			}).Errorf("unable to retrieve activity detail: \n%v", err)
			res.Render(500, map[string]interface{}{
				"error":    "Unable to retrieve activity detail",
				"activity": activityDetail,
			})
			return
		}

		// range over segment efforts from the activity detail
		// to obtain segment details to cache
		for _, effort := range activityDetail.SegmentEfforts {
			logrus.WithField("SEGMENT", effort.Segment.Name).Info("segment effort from activity detail")
			// check if segment is in MongoDB
			segment, err := models.GetSegmentByID(effort.Segment.Id)
			if err != nil {
				// segment not found, make request to strava
				logrus.WithField("SEGMENT ID", effort.Segment.Id).Infof("segment %v not found in MongoDB... saving", effort.Segment.Id)
				segmentDetail, err := strava.NewSegmentsService(client).Get(effort.Segment.Id).Do()
				if err != nil {
					logrus.WithFields(map[string]interface{}{
						"SEGMENT NAME": effort.Segment.Name,
						"SEGMENT ID":   effort.Segment.Id,
					}).Errorf("unable to retrieve segment detail for %d %s", effort.Segment.Id, effort.Segment.Name)
					res.Render(500, map[string]interface{}{
						"error":   "Unable to retrieve segment detail",
						"segment": effort.Segment,
					})
					return
				}
				logrus.Infof("rate limit percent: %v", strava.RateLimiting.FractionReached()*100)
				logrus.WithField("SEGMENT DETAIL ID", segmentDetail.Id).Infof("segment %d returned from strava", segmentDetail.Id)
				segment, err := models.SaveSegment(segmentDetail)
				if err != nil {
					logrus.WithError(err).Errorf("unable to save segment detail %d to MongoDB", segmentDetail.Id)
					return
				}
				logrus.WithField("SEGMENT ID", segment.ID).Infof("retrieved segment %d detail from strava", segment.ID)
				segmentMap[segment.ID] = segment
			} else {
				// segment was found and returned
				segmentMap[segment.ID] = segment
			}
		}
	}
	// store segment map for user

	// range over segmentMap to creake segments array
	for _, segment := range segmentMap {
		segments = append(segments, segment)
	}
	logrus.Infof("found %d unique segments", len(segments))
	res.Render(200, segments)
}
