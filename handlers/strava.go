package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/jrzimmerman/bestrida-server-go/models"
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
	logrus.WithField("USER TOKEN", userToken).Info("user token on oAuth success")
	logrus.WithField("USER ID", userID).Info("user id on oAuth success")
	url := "/login.html?oauth_token=" + userToken + "&userId=" + userID
	http.Redirect(w, r, url, 301)

	_, err := models.RegisterUser(auth)
	if err != nil {
		logrus.Errorf("error registering user %d", auth.Athlete.Id)
	}
}

func oAuthFailure(err error, w http.ResponseWriter, r *http.Request) {
	if err == strava.OAuthAuthorizationDeniedErr {
		logrus.WithError(err).Error("The user clicked the 'Do not Authorize' button on the previous page.")
	} else if err == strava.OAuthInvalidCredentialsErr {
		logrus.WithError(err).Error("You provided an incorrect client_id or client_secret.")
	} else if err == strava.OAuthInvalidCodeErr {
		logrus.WithError(err).Error("The temporary token was not recognized, this shouldn't happen normally")
	} else if err == strava.OAuthServerErr {
		logrus.WithError(err).Error("There was some sort of server error, try again to see if the problem continues")
	} else {
		logrus.WithError(err).Error("authorization failure")
	}
}
