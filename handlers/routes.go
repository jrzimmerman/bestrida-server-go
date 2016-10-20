package handlers

import (
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
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
		s.GET("athletes/:id/segments", GetSegmentsByUserIDFromStrava)
	}

	api := r.Group("/api")
	{
		api.GET("users/:id", GetUserByID)
		api.GET("segments/:id", GetSegmentByID)
		api.GET("challenges/:id", GetChallengeByID)
	}

	return r
}
