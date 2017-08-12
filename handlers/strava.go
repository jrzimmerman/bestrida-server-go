package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jrzimmerman/bestrida-server-go/models"
	log "github.com/sirupsen/logrus"
	"github.com/strava/go.strava"
)

// AuthHandler Route to display Auth button
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	// you should make this a template in your real application
	fmt.Fprintf(w, `<a href="%s">`, authenticator.AuthorizationURL("", strava.Permissions.Public, false))
	fmt.Fprint(w, `<img src="/assets/img/btn_strava_connectwith_orange.png" />`)
	fmt.Fprint(w, `</a>`)
}

func oAuthSuccess(auth *strava.AuthorizationResponse, w http.ResponseWriter, r *http.Request) {
	userToken := string(auth.AccessToken)
	userID := strconv.Itoa(int(auth.Athlete.Id))
	log.WithField("USER TOKEN", userToken).Info("user token on oAuth success")
	log.WithField("USER ID", userID).Info("user id on oAuth success")
	url := "/login.html?oauth_token=" + userToken + "&userId=" + userID
	http.Redirect(w, r, url, http.StatusFound)

	_, err := models.RegisterUser(auth)
	if err != nil {
		log.Errorf("error registering user %d", auth.Athlete.Id)
	}
}

func oAuthFailure(err error, w http.ResponseWriter, r *http.Request) {
	if err == strava.OAuthAuthorizationDeniedErr {
		log.WithError(err).Error("The user clicked the 'Do not Authorize' button on the previous page.")
	} else if err == strava.OAuthInvalidCredentialsErr {
		log.WithError(err).Error("You provided an incorrect client_id or client_secret.")
	} else if err == strava.OAuthInvalidCodeErr {
		log.WithError(err).Error("The temporary token was not recognized, this shouldn't happen normally")
	} else if err == strava.OAuthServerErr {
		log.WithError(err).Error("There was some sort of server error, try again to see if the problem continues")
	} else {
		log.WithError(err).Error("authorization failure")
	}
}
