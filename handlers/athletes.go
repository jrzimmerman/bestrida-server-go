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
		log.Error("unable to convert ID param")
		res.Render(http.StatusBadRequest, map[string]interface{}{
			"error": "unable to convert ID param",
			"stack": err,
		})
		return
	}

	user, err := models.GetUserByID(numID)
	if err != nil {
		log.WithField("ID", numID).Error("unable to retrieve user from database")
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "unable to retrieve user from database",
			"stack": err,
		})
		return
	}

	client := strava.NewClient(user.Token)

	log.Info("Fetching athlete info...\n")
	// retrieve a list of users segments from Strava API
	athlete, err := strava.NewCurrentAthleteService(client).Get().Do()
	if err != nil {
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "Unable to retrieve athlete info",
			"stack": err,
		})
		return
	}
	log.Infof("rate limit percent: %v", strava.RateLimiting.FractionReached()*100)
	log.Infof("athlete %v retrieved from strava", athlete.Id)

	u, err := user.UpdateAthlete(athlete)
	if err != nil {
		log.WithError(err).Errorf("unable to update athlete %d", athlete.Id)
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "Unable to update athlete in DB",
			"stack": err,
		})
		return
	}
	res.Render(http.StatusOK, &u)
}

// GetFriendsFromStrava gets a users friends from Strava
func GetFriendsFromStrava(numID int64) (friends []*models.Friend, err error) {
	// get user based on ID
	user, err := models.GetUserByID(numID)
	if err != nil {
		log.WithField("USER ID", numID).Errorf("unable to retrieve user %v from database", numID)
		return nil, err
	}

	// friends map to store unique friend info
	friendMap := make(map[int64]*models.Friend, len(user.Friends))
	// iterate over users existing friends to populate friend map
	for _, friend := range user.Friends {
		friendMap[friend.ID] = friend
	}

	client := strava.NewClient(user.Token)

	log.Info("Fetching athlete friends info...\n")
	// retrieve a list of users friends from Strava API
	stravaFriends, err := strava.NewCurrentAthleteService(client).ListFriends().Do()
	if err != nil {
		return nil, err
	}
	log.Infof("rate limit percent: %v", strava.RateLimiting.FractionReached()*100)
	log.Infof("Finished fetching %v athlete friends from Strava...\n", len(stravaFriends))

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
			count := 0
			wins := 0
			losses := 0
			if friend, ok := friendMap[stravaFriend.Id]; ok {
				count = friend.ChallengeCount
				wins = friend.Wins
				losses = friend.Losses
			}
			// existing friend, update info and keep count
			friendMap[stravaFriend.Id] = &models.Friend{
				ID:             stravaFriend.Id,
				FirstName:      stravaFriend.FirstName,
				LastName:       stravaFriend.LastName,
				FullName:       stravaFriend.FirstName + " " + stravaFriend.LastName,
				Photo:          stravaFriend.Profile,
				ChallengeCount: count,
				Wins:           wins,
				Losses:         losses,
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
		return nil, err
	}
	return friends, nil
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
			"stack": err,
		})
		return
	}

	// Get users friends from strava and save to DB
	friends, err := GetFriendsFromStrava(numID)
	if err != nil {
		log.WithField("USER ID", numID).Errorf("unable to retrieve user %v friends from strava", numID)
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "unable to retrieve users friends from strava",
			"stack": err,
		})
		return
	}

	log.Infof("found %d friends for athlete %d from strava", len(friends), numID)
	res.Render(http.StatusOK, friends)
}

// GetUserSegmentsFromStrava gets a users recently completed segments from Strava
func GetUserSegmentsFromStrava(numID int64, page int) (userSegmentSlice []*models.UserSegment, err error) {
	// find user by numID to retrieve strava token
	user, err := models.GetUserByID(numID)
	if err != nil {
		log.WithField("ID", numID).Error("unable to retrieve user from database")
		return nil, err
	}

	// unique segment map to store a users segments
	userSegments := make(map[int64]*models.UserSegment, len(user.Segments))
	for _, segment := range user.Segments {
		userSegments[segment.ID] = segment
	}

	// create new strava client with user token
	client := strava.NewClient(user.Token)

	log.Info("Fetching starred segments from Strava...\n")
	starred, err := strava.NewCurrentAthleteService(client).ListStarredSegments().Do()
	if err != nil {
		return nil, err
	}

	// range over a users starred segments
	// to obtain segment details to cache
	for _, seg := range starred {
		log.Infof("segment %v was starred by user", seg.Id)
		log.WithField("SEGMENT", seg.Name).Info("segment effort from activity detail")
		// check if segment is in database
		segment, err := models.GetSegmentByID(seg.Id)
		if err != nil {
			// segment not found, make request to strava
			log.WithField("SEGMENT ID", seg.Id).Infof("segment %v not found in database... saving", seg.Id)
			segmentDetail, err := strava.NewSegmentsService(client).Get(seg.Id).Do()
			if err != nil {
				log.WithFields(log.Fields{
					"SEGMENT NAME": seg.Name,
					"SEGMENT ID":   seg.Id,
				}).Errorf("unable to retrieve segment detail for %d %s", seg.Id, seg.Name)
				return nil, err
			}
			log.Infof("rate limit percent: %v", strava.RateLimiting.FractionReached()*100)
			log.WithField("SEGMENT DETAIL ID", segmentDetail.Id).Infof("segment %d returned from strava", segmentDetail.Id)
			saved, err := models.SaveSegment(segmentDetail)
			if err != nil {
				log.WithError(err).Errorf("unable to save segment detail %d to database", segmentDetail.Id)
				return nil, err
			}
			log.WithField("SEGMENT ID", saved.ID).Infof("segment %d stored in DB", saved.ID)

			// store saved segment in userSegments map
			userSegments[saved.ID] = &models.UserSegment{
				ID:           saved.ID,
				Name:         saved.Name,
				ActivityType: saved.ActivityType,
				Count:        0,
			}

		} else if time.Now().After(segment.UpdatedAt.AddDate(0, 0, 7)) {
			// update segment in DB if segment data is stale
			// this is required by Strava API license agreement
			log.WithField("SEGMENT ID", seg.Id).Infof("segment %v is greater than 7 days old... updating", seg.Id)
			segmentDetail, err := strava.NewSegmentsService(client).Get(seg.Id).Do()
			if err != nil {
				log.WithFields(log.Fields{
					"SEGMENT NAME": seg.Name,
					"SEGMENT ID":   seg.Id,
				}).Errorf("unable to retrieve segment detail for %d %s", seg.Id, seg.Name)
				return nil, err
			}
			log.Infof("rate limit percent: %v", strava.RateLimiting.FractionReached()*100)
			log.WithField("SEGMENT DETAIL ID", segmentDetail.Id).Infof("segment %d returned from strava", segmentDetail.Id)
			updated, err := segment.UpdateSegment(segmentDetail)
			if err != nil {
				log.WithError(err).Errorf("unable to save segment detail %d to database", segmentDetail.Id)
				return nil, err
			}
			log.WithField("SEGMENT ID", updated.ID).Infof("segment %d updated in DB", updated.ID)

			count := 0
			if segment, ok := userSegments[segment.ID]; ok {
				count = segment.Count
			}
			// store updated segment in userSegments map
			userSegments[updated.ID] = &models.UserSegment{
				ID:           updated.ID,
				Name:         updated.Name,
				ActivityType: updated.ActivityType,
				Count:        count,
			}
		} else {
			count := 0
			if segment, ok := userSegments[segment.ID]; ok {
				count = segment.Count
			}
			// segment was found and returned
			// add segment to userSegments
			userSegments[segment.ID] = &models.UserSegment{
				ID:           segment.ID,
				Name:         segment.Name,
				ActivityType: segment.ActivityType,
				Count:        count,
			}
		}
	}

	log.Info("Fetching athlete activities from Strava...\n")
	activities, err := strava.NewCurrentAthleteService(client).ListActivities().Page(1).PerPage(page).Do()
	if err != nil {
		return nil, err
	}
	log.Infof("rate limit percent: %v", strava.RateLimiting.FractionReached()*100)
	log.Infof("Finished fetching %v athlete activities from Strava...\n", len(activities))
	// range over activity summary to get activity details
	// the activity summary does not contain segment effort information
	for _, activitySummary := range activities {
		log.WithFields(log.Fields{
			"NAME": activitySummary.Name,
			"ID":   activitySummary.Id,
		}).Info("activity summary")

		// request activity detail from strava to obtain segment effort information
		activityDetail, err := strava.NewActivitiesService(client).Get(activitySummary.Id).IncludeAllEfforts().Do()
		if err != nil {
			log.WithFields(log.Fields{
				"NAME": activityDetail.Name,
				"ID":   activityDetail.Id,
			}).Errorf("unable to retrieve activity detail: \n%v", err)
			return nil, err
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
					return nil, err
				}
				log.Infof("rate limit percent: %v", strava.RateLimiting.FractionReached()*100)
				log.WithField("SEGMENT DETAIL ID", segmentDetail.Id).Infof("segment %d returned from strava", segmentDetail.Id)
				saved, err := models.SaveSegment(segmentDetail)
				if err != nil {
					log.WithError(err).Errorf("unable to save segment detail %d to database", segmentDetail.Id)
					return nil, err
				}
				log.WithField("SEGMENT ID", saved.ID).Infof("segment %d stored in DB", saved.ID)

				// store saved segment in userSegments map
				userSegments[saved.ID] = &models.UserSegment{
					ID:           saved.ID,
					Name:         saved.Name,
					ActivityType: saved.ActivityType,
					Count:        0,
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
					return nil, err
				}
				log.Infof("rate limit percent: %v", strava.RateLimiting.FractionReached()*100)
				log.WithField("SEGMENT DETAIL ID", segmentDetail.Id).Infof("segment %d returned from strava", segmentDetail.Id)
				updated, err := segment.UpdateSegment(segmentDetail)
				if err != nil {
					log.WithError(err).Errorf("unable to save segment detail %d to database", segmentDetail.Id)
					return nil, err
				}
				log.WithField("SEGMENT ID", updated.ID).Infof("segment %d updated in DB", updated.ID)

				count := 0
				if segment, ok := userSegments[segment.ID]; ok {
					count = segment.Count
				}
				// store updated segment in userSegments map
				userSegments[updated.ID] = &models.UserSegment{
					ID:           updated.ID,
					Name:         updated.Name,
					ActivityType: updated.ActivityType,
					Count:        count,
				}
			} else {
				count := 0
				if segment, ok := userSegments[segment.ID]; ok {
					count = segment.Count
				}
				// segment was found and returned
				// add segment to userSegments
				userSegments[segment.ID] = &models.UserSegment{
					ID:           segment.ID,
					Name:         segment.Name,
					ActivityType: segment.ActivityType,
					Count:        count,
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
		return nil, err
	}
	return userSegmentSlice, nil
}

// GetSegmentsByUserIDFromStrava returns a list of segments for a specific user by ID from strava
func GetSegmentsByUserIDFromStrava(w http.ResponseWriter, r *http.Request) {
	res := New(w)
	id := chi.URLParam(r, "id")

	// convert user id string from url param to number
	numID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.WithField("USER ID", numID).Error("unable to convert USER ID param")
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "unable to convert USER ID param",
			"stack": err,
		})
		return
	}

	// page is the amount of activites to request from Strava
	var page = 30
	// get users segments from Strava
	userSegments, err := GetUserSegmentsFromStrava(numID, page)
	if err != nil {
		log.WithField("USER ID", numID).Errorf("unable to retrieve user %v segments from strava", numID)
		res.Render(http.StatusInternalServerError, map[string]interface{}{
			"error": "unable to retrieve users segments from strava",
			"stack": err,
		})
		return
	}

	log.WithField("USER ID", numID).Infof("found %d segments for user %v", len(userSegments), numID)
	res.Render(http.StatusOK, userSegments)
}
