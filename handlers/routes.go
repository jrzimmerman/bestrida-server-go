package handlers

import (
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	strava "github.com/strava/go.strava"
	"gopkg.in/gin-gonic/gin.v1"
)

var authenticator *strava.OAuthAuthenticator
var clientID = getEnvString("STRAVA_CLIENT_ID")
var clientSecret = getEnvString("STRAVA_CLIENT_SECRET")
var accessToken = getEnvString("STRAVA_ACCESS_TOKEN")
var port = getEnvString("PORT")

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

func getEnvString(env string) string {
	str, ok := os.LookupEnv(env)
	if !ok {
		log.WithField("ENV", env).Fatal("Missing required environment variable")
	}
	return str
}
