package handlers

import (
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"

	"github.com/jrzimmerman/bestrida-server-go/models"
)

func GetSingleAthlete(c *gin.Context) {
	id := 123

	user, err := models.GetUserByID(id)
	if err != nil {
		logrus.WithField("id", id).Debug("unable to get user by ID")
	}

	c.JSON(200, user)
}
