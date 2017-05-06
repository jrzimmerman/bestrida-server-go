package handlers

import (
	"net/http"

	"github.com/jrzimmerman/bestrida-server-go/utils"
	strava "github.com/strava/go.strava"
	"gopkg.in/gin-gonic/gin.v1"
)

var authenticator *strava.OAuthAuthenticator
var clientID = utils.GetEnvString("STRAVA_CLIENT_ID")
var clientSecret = utils.GetEnvString("STRAVA_CLIENT_SECRET")
var accessToken = utils.GetEnvString("STRAVA_ACCESS_TOKEN")
var port = utils.GetEnvString("PORT")

// API initializes routes with Gin
func API() http.Handler {
	r := gin.Default()
	r.Use(CORS())
	s := r.Group("/strava")
	{
		s.GET("auth", StravaAuth)
		s.GET("athletes/:id", GetAthleteByIDFromStrava)
		s.GET("athletes/:id/friends", GetFriendsByUserIDFromStrava)
		s.GET("athletes/:id/segments", GetSegmentsByUserIDFromStrava)
		s.GET("athletes/:id/segments/:segmentID/efforts", GetEffortsBySegmentIDFromStrava)
		s.GET("segments/:id", GetSegmentByIDFromStrava)
	}

	api := r.Group("/api")
	{
		api.GET("users/:id", GetUserByID)
		api.GET("segments/:id", GetSegmentByID)
		api.GET("challenges/:id", GetChallengeByID)
	}

	return r
}
