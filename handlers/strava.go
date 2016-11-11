package handlers

import "gopkg.in/gin-gonic/gin.v1"

// StravaAuth will render a request to authorize strava user with Bestrida
func StravaAuth(c *gin.Context) {
	c.JSON(200, map[string]interface{}{
		"ID":     clientID,
		"SECRET": clientSecret,
	})
}
