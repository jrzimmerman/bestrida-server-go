package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jrzimmerman/bestrida-server-go/handlers"
)

// API initializes routes with Gin
func API() http.Handler {
	r := gin.Default()
	api := r.Group("/api")
	{
		api.GET("athletes/:id", handlers.GetSingleAthlete)
	}
	return r
}
