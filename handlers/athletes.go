package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/jrzimmerman/bestrida-server-go/models"
	log "github.com/sirupsen/logrus"
	strava "github.com/strava/go.strava"
)

// GetAthleteByIDFromStrava returns the strava athlete with the specified ID
func GetAthleteByIDFromStrava(w http.ResponseWriter, r *http.Request) {
	res := New(w)

	id := chi.URLParam(r, "id")

	numID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.WithField("ID", numID).Error("unable to convert ID param")
		res.Render(http.StatusBadRequest, "unable to convert ID param")
		return
	}

	user, err := models.GetUserByID(numID)
	if err != nil {
		log.WithField("ID", numID).Error("unable to retrieve user from database")
		res.Render(http.StatusInternalServerError, map[string]interface{}{"error": "unable to retrieve user from database"})
		return
	}

	client := strava.NewClient(user.Token)

	log.Info("Fetching athlete info...\n")
	// retrieve a list of users segments from Strava API
	athlete, err := strava.NewCurrentAthleteService(client).Get().Do()
	if err != nil {
		res.Render(http.StatusInternalServerError, map[string]interface{}{"error": "Unable to retrieve athlete info"})
		return
	}

	log.Infof("rate limit percent: %v", strava.RateLimiting.FractionReached()*100)
	log.Infof("athlete %v retrieved from strava", athlete.Id)
	u, err := user.UpdateAthlete(athlete)
	if err != nil {
		log.WithError(err).Errorf("unable to update athlete %d", athlete.Id)
	}
	res.Render(http.StatusOK, &u)
}

// GetFriendsByUserIDFromStrava returns a list of friends for a specific user by ID from strava
func GetFriendsByUserIDFromStrava(w http.ResponseWriter, r *http.Request) {
	res := New(w)

	id := chi.URLParam(r, "id")

	numID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.WithField("ID", numID).Error("unable to convert ID param")
		res.Render(http.StatusInternalServerError, map[string]interface{}{"error": "unable to convert ID param"})
		return
	}

	user, err := models.GetUserByID(numID)
	if err != nil {
		log.WithField("USER ID", numID).Errorf("unable to retrieve user %v from database", numID)
		res.Render(http.StatusInternalServerError, map[string]interface{}{"error": "unable to retrieve user from database"})
		return
	}

	client := strava.NewClient(user.Token)

	log.Info("Fetching athlete friends info...\n")
	// retrieve a list of users friends from Strava API
	friends, err := strava.NewCurrentAthleteService(client).ListFriends().Do()
	if err != nil {
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "Unable to retrieve athlete friends",
		})
		return
	}

	log.Infof("%v friends from strava", len(friends))

	for _, friend := range friends {
		log.Info("friend ID : \n %v", &friend.Id)
	}
	// store friends for user
	// err = user.SaveUserFriends(friends)
	// if err != nil {
	// 	log.WithError(err).Errorf("unable to save user friends for user %d to database", user.ID)
	// 	return
	// }

	// return request
	log.Infof("rate limit percent: %v", strava.RateLimiting.FractionReached()*100)
	log.Infof("found %d friends for athlete %d from strava", len(friends), user.ID)
	res.Render(http.StatusOK, friends)
}

// GetSegmentsByUserIDFromStrava returns a list of segments for a specific user by ID from strava
func GetSegmentsByUserIDFromStrava(w http.ResponseWriter, r *http.Request) {
	res := New(w)

	id := chi.URLParam(r, "id")

	var userSegmentSlice []*models.UserSegment

	// convert user id string from url param to number
	numID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.WithField("USER ID", numID).Error("unable to convert USER ID param")
		res.Render(http.StatusInternalServerError, map[string]interface{}{"error": "unable to convert USER ID param"})
		return
	}

	// find user by numID to retrieve strava token
	user, err := models.GetUserByID(numID)
	if err != nil {
		log.WithField("ID", numID).Error("unable to retrieve user from database")
		res.Render(404, "unable to retrieve user from database")
		return
	}

	userSegments := make(map[int64]*models.UserSegment, 0)
	for _, segment := range user.Segments {
		userSegments[segment.ID] = segment
	}

	// create new strava client with user token
	client := strava.NewClient(user.Token)

	log.Info("Fetching athlete activity summary info...\n")
	activities, err := strava.NewCurrentAthleteService(client).ListActivities().Page(1).PerPage(200).Do()
	if err != nil {
		res.Render(http.StatusInternalServerError, map[string]interface{}{"error": "Unable to retrieve athlete activities summary"})
		return
	}

	// create a segment map to return only unique segments
	segmentMap := make(map[int64]models.Segment, 0)

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
			}).Errorf("unable to retrieve activity detail: \n%v", err)
			res.Render(http.StatusInternalServerError, map[string]interface{}{
				"error":    "Unable to retrieve activity detail",
				"activity": activityDetail,
			})
			return
		}

		// range over segment efforts from the activity detail
		// to obtain segment details to cache
		for _, effort := range activityDetail.SegmentEfforts {
			log.WithField("SEGMENT", effort.Segment.Name).Info("segment effort from activity detail")
			// check if segment is in database
			segment, err := models.GetSegmentByID(effort.Segment.Id)
			if err != nil {
				// segment not found, make request to strava
				log.WithField("SEGMENT ID", effort.Segment.Id).Infof("segment %v not found in database... saving", effort.Segment.Id)
				segmentDetail, err := strava.NewSegmentsService(client).Get(effort.Segment.Id).Do()
				if err != nil {
					log.WithFields(log.Fields{
						"SEGMENT NAME": effort.Segment.Name,
						"SEGMENT ID":   effort.Segment.Id,
					}).Errorf("unable to retrieve segment detail for %d %s", effort.Segment.Id, effort.Segment.Name)
					res.Render(http.StatusInternalServerError, map[string]interface{}{
						"error":   "Unable to retrieve segment detail",
						"segment": effort.Segment,
					})
					return
				}
				log.Infof("rate limit percent: %v", strava.RateLimiting.FractionReached()*100)
				log.WithField("SEGMENT DETAIL ID", segmentDetail.Id).Infof("segment %d returned from strava", segmentDetail.Id)
				segment, err := models.SaveSegment(segmentDetail)
				if err != nil {
					log.WithError(err).Errorf("unable to save segment detail %d to database", segmentDetail.Id)
					return
				}
				log.WithField("SEGMENT ID", segment.ID).Infof("retrieved segment %d detail from strava", segment.ID)
				// add segment to segmentMap if not found
				if _, ok := segmentMap[segment.ID]; !ok {
					segmentMap[segment.ID] = *segment
				}
				// add segment to userSegments if not found
				if _, ok := userSegments[segment.ID]; !ok {
					userSegments[segment.ID] = &models.UserSegment{
						ID:    segment.ID,
						Name:  segment.Name,
						Count: 0,
					}
				}
			} else {
				// segment was found and returned
				// add segment to segmentMap if not found
				if _, ok := segmentMap[segment.ID]; !ok {
					segmentMap[segment.ID] = *segment
				}
				// add segment to userSegments if not found
				if _, ok := userSegments[segment.ID]; !ok {
					userSegments[segment.ID] = &models.UserSegment{
						ID:    segment.ID,
						Name:  segment.Name,
						Count: 0,
					}
				}
			}
		}
	}

	for _, userSegment := range userSegments {
		userSegmentSlice = append(userSegmentSlice, userSegment)
	}

	// store segment map for user
	err = user.SaveUserSegments(userSegmentSlice)
	if err != nil {
		log.WithError(err).Errorf("unable to save user segments for user %d to database", user.ID)
		return
	}

	log.WithField("USER ID", user.ID).Infof("found %d segments for user %v", len(userSegmentSlice), user.ID)
	res.Render(http.StatusOK, userSegmentSlice)
}
