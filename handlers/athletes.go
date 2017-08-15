package handlers

import (
	"net/http"
	"strconv"
	"time"

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
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "unable to convert ID param",
		})
		return
	}

	// get user based on ID
	user, err := models.GetUserByID(numID)
	if err != nil {
		log.WithField("USER ID", numID).Errorf("unable to retrieve user %v from database", numID)
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "unable to retrieve user from database",
		})
		return
	}

	// friends slice to store in db and return in response
	var friends []*models.Friend
	// friends map to store unique friend info
	friendMap := make(map[int64]*models.Friend, 0)
	// iterate over users existing friends to populate friend map
	for _, friend := range user.Friends {
		friendMap[friend.ID] = friend
	}

	client := strava.NewClient(user.Token)

	log.Info("Fetching athlete friends info...\n")
	// retrieve a list of users friends from Strava API
	stravaFriends, err := strava.NewCurrentAthleteService(client).ListFriends().Do()
	if err != nil {
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "Unable to retrieve athlete friends",
		})
		return
	}
	log.Infof("rate limit percent: %v", strava.RateLimiting.FractionReached()*100)
	log.Infof("%v friends from strava", len(stravaFriends))

	// update friends map based upon strava friend data
	for _, stravaFriend := range stravaFriends {
		if _, ok := friendMap[stravaFriend.Id]; !ok {
			// new friend set new information
			friendMap[stravaFriend.Id] = &models.Friend{
				ID:        stravaFriend.Id,
				FirstName: stravaFriend.FirstName,
				LastName:  stravaFriend.LastName,
				FullName:  stravaFriend.FirstName + " " + stravaFriend.LastName,
				Photo:     stravaFriend.Profile,
			}
		} else {
			// existing friend, update info and keep count
			friendMap[stravaFriend.Id] = &models.Friend{
				ID:             stravaFriend.Id,
				FirstName:      stravaFriend.FirstName,
				LastName:       stravaFriend.LastName,
				FullName:       stravaFriend.FirstName + " " + stravaFriend.LastName,
				Photo:          stravaFriend.Profile,
				ChallengeCount: friendMap[stravaFriend.Id].ChallengeCount,
				Wins:           friendMap[stravaFriend.Id].Wins,
				Losses:         friendMap[stravaFriend.Id].Losses,
			}
		}
	}

	// range over unique friends to store updated info in slice
	for _, friend := range friendMap {
		friends = append(friends, friend)
	}

	// store friends slice for user
	err = user.SaveUserFriends(friends)
	if err != nil {
		log.WithError(err).Errorf("unable to save user friends for user %d to database", user.ID)
		return
	}

	log.Infof("found %d friends for athlete %d from strava", len(friends), user.ID)
	res.Render(http.StatusOK, friends)
}

// GetSegmentsByUserIDFromStrava returns a list of segments for a specific user by ID from strava
func GetSegmentsByUserIDFromStrava(w http.ResponseWriter, r *http.Request) {
	res := New(w)
	id := chi.URLParam(r, "id")

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

	// unique segment slice to store a users segments
	var userSegmentSlice []*models.UserSegment

	// unique segment map to store a users segments
	userSegments := make(map[int64]*models.UserSegment, 0)
	for _, segment := range user.Segments {
		userSegments[segment.ID] = segment
	}

	// create new strava client with user token
	client := strava.NewClient(user.Token)

	log.Info("Fetching athlete activity summary info from Strava...\n")
	activities, err := strava.NewCurrentAthleteService(client).ListActivities().Page(1).PerPage(200).Do()
	if err != nil {
		res.Render(http.StatusInternalServerError, map[string]interface{}{"error": "Unable to retrieve athlete activities summary"})
		return
	}
	log.Infof("rate limit percent: %v", strava.RateLimiting.FractionReached()*100)

	// range over activity summary to get activity details
	// the activity summary does not contain segment effort information
	for _, activitySummary := range activities {
		log.WithFields(log.Fields{
			"NAME": activitySummary.Name,
			"ID":   activitySummary.Id,
		}).Info("activity summary")

		// request activity detail from strava to obtain segment effort information
		activityDetail, err := strava.NewActivitiesService(client).Get(activitySummary.Id).Do()
		if err != nil {
			log.WithFields(log.Fields{
				"NAME": activityDetail.Name,
				"ID":   activityDetail.Id,
			}).Errorf("unable to retrieve activity detail: \n%v", err)
			res.Render(http.StatusInternalServerError, map[string]interface{}{
				"error":    "Unable to retrieve activity detail",
				"activity": activityDetail,
			})
			return
		}
		log.Infof("rate limit percent: %v", strava.RateLimiting.FractionReached()*100)

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
				saved, err := models.SaveSegment(segmentDetail)
				if err != nil {
					log.WithError(err).Errorf("unable to save segment detail %d to database", segmentDetail.Id)
					return
				}
				log.WithField("SEGMENT ID", saved.ID).Infof("segment %d stored in DB", saved.ID)

				// store saved segment in userSegments map
				userSegments[saved.ID] = &models.UserSegment{
					ID:           saved.ID,
					Name:         saved.Name,
					ActivityType: saved.ActivityType,
				}

			} else if time.Now().After(segment.UpdatedAt.AddDate(0, 0, 7)) {
				// update segment in DB if segment data is stale
				// this is required by Strava API license agreement
				log.WithField("SEGMENT ID", effort.Segment.Id).Infof("segment %v is greater than 7 days old... updating", effort.Segment.Id)
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
				updated, err := segment.UpdateSegment(segmentDetail)
				if err != nil {
					log.WithError(err).Errorf("unable to save segment detail %d to database", segmentDetail.Id)
					return
				}
				log.WithField("SEGMENT ID", updated.ID).Infof("segment %d updated in DB", updated.ID)

				// store updated segment in userSegments map
				userSegments[updated.ID] = &models.UserSegment{
					ID:           updated.ID,
					Name:         updated.Name,
					ActivityType: updated.ActivityType,
					Count:        userSegments[updated.ID].Count,
				}
			} else {
				// segment was found and returned
				// add segment to userSegments
				userSegments[segment.ID] = &models.UserSegment{
					ID:           segment.ID,
					Name:         segment.Name,
					ActivityType: segment.ActivityType,
					Count:        userSegments[segment.ID].Count,
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

	log.WithField("USER ID", user.ID).Infof("updated %d segments for user %v", len(userSegmentSlice), user.ID)
	res.Render(http.StatusOK, userSegmentSlice)
}
