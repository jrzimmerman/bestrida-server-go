package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jrzimmerman/bestrida-server-go/handlers"
	"github.com/jrzimmerman/bestrida-server-go/middleware"
)

// API initializes routes with Gin
func API() http.Handler {
	r := gin.Default()
	r.Use(middleware.CORS())
	auth := r.Group("/auth")
	{
		auth.GET("strava", handlers.StravaAuth)
	}

	api := r.Group("/api")
	{
		api.GET("users/:id", handlers.GetUserByID)
		api.GET("segments/:id", handlers.GetSegmentByID)
		api.GET("challenges/:id", handlers.GetChallengeByID)
	}

	return r
}
