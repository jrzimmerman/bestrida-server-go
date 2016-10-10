package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// API initializes routes with Gin
func API() http.Handler {
	r := gin.Default()
	r.Use(CORS())
	s := r.Group("/strava")
	{
		s.GET("auth", StravaAuth)
		s.GET("athletes/:id", GetAthleteByIDFromStrava)
		s.GET("athletes/:id/friends", GetFriendsByUserIDFromStrava)
	}

	api := r.Group("/api")
	{
		api.GET("users/:id", GetUserByID)
		api.GET("segments/:id", GetSegmentByID)
		api.GET("challenges/:id", GetChallengeByID)
	}

	return r
}
