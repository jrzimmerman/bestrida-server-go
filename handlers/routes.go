package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// API initializes routes with Gin
func API() http.Handler {
	r := gin.Default()
	r.Use(CORS())
	auth := r.Group("/auth")
	{
		auth.GET("strava", StravaAuth)
	}

	api := r.Group("/api")
	{
		api.GET("users/:id", GetUserByID)
		api.GET("segments/:id", GetSegmentByID)
		api.GET("challenges/:id", GetChallengeByID)
	}

	return r
}
