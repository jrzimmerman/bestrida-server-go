package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	strava "github.com/strava/go.strava"
	"gopkg.in/gin-gonic/gin.v1"
)

var authenticator *strava.OAuthAuthenticator
var clientID = getEnvString("STRAVA_CLIENT_ID")
var clientSecret = getEnvString("STRAVA_CLIENT_SECRET")

// StravaAuth will render a request to authorize strava user with Bestrida
func StravaAuth(c *gin.Context) {

	// TODO: send auth request to Strava
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
		log.WithError(err).Error("This is the main error your application should handle.")
	} else if err == strava.OAuthInvalidCredentialsErr {
		log.WithError(err).Error("You provided an incorrect client_id or client_secret.")
		log.WithError(err).Error("Did you remember to set them at the beginning of this file?")
	} else if err == strava.OAuthInvalidCodeErr {
		log.WithError(err).Error("The temporary token was not recognized, this shouldn't happen normally")
	} else if err == strava.OAuthServerErr {
		log.WithError(err).Error("There was some sort of server error, try again to see if the problem continues")
	} else {
		log.WithError(err).Error("Unknown Error Authorizing")
	}
}

func getEnvString(env string) string {
	str, ok := os.LookupEnv(env)
	if !ok {
		log.WithField("ENV", env).Fatal("Missing required environment variable")
	}
	return str
}
