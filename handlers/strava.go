package handlers

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	strava "github.com/strava/go.strava"
	"gopkg.in/gin-gonic/gin.v1"
)

// StravaAuth will render a request to authorize strava user with Bestrida
func StravaAuth(c *gin.Context) {
	c.JSON(200, map[string]interface{}{
		"ID":     clientID,
		"SECRET": clientSecret,
	})
}

func oAuthSuccess(auth *strava.AuthorizationResponse, w http.ResponseWriter, r *http.Request) {
	log.Info("SUCCESS: At this point you can use this information to create a new user or link the account to one of your existing users")
	log.WithField("State", auth.State).Info("State Returned From Strava")
	log.WithField("Access Token", auth.AccessToken).Info("Strava Access Token")

	log.Info("The Authenticated Athlete (you): ")
	content, _ := json.MarshalIndent(auth.Athlete, "", " ")
	log.WithField("Content", content).Info("Content From Strava")
}

func oAuthFailure(err error, w http.ResponseWriter, r *http.Request) {
	log.Error("Authorization Failure: ")

	// some standard error checking
	if err == strava.OAuthAuthorizationDeniedErr {
		log.WithError(err).Error("The user clicked the 'Do not Authorize' button on the previous page.")
	} else {
		log.WithError(err).Error("Unknown Error Authorizing with Strava")
	}
}
