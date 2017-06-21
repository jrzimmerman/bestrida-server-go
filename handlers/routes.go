package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/jrzimmerman/bestrida-server-go/utils"
	"github.com/pressly/chi"
	"github.com/strava/go.strava"
)

var authenticator = &strava.OAuthAuthenticator{
	CallbackURL:            fmt.Sprintf("http://localhost:%s/strava/auth/callback", port),
	RequestClientGenerator: nil,
}
var clientID = utils.GetEnvString("STRAVA_CLIENT_ID")
var clientSecret = utils.GetEnvString("STRAVA_CLIENT_SECRET")
var accessToken = utils.GetEnvString("STRAVA_ACCESS_TOKEN")
var port = utils.GetEnvString("PORT")

// API initializes all endpoints
func API() (mux *chi.Mux) {
	path, err := authenticator.CallbackPath()
	if err != nil {
		// possibly that the callback url set above is invalid
		logrus.Error(err)
	}
	clientIDInt, err := strconv.Atoi(clientID)
	if err != nil {
		logrus.Error(err)
	}
	strava.ClientId = clientIDInt
	strava.ClientSecret = clientSecret

	mux = chi.NewRouter()
	mux.Use(CORS)

	mux.HandleFunc(path, authenticator.HandlerFunc(oAuthSuccess, oAuthFailure))
	mux.FileServer("/", http.Dir("public"))

	mux.Route("/api", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Route("/:id", func(r chi.Router) {
				r.Get("/", GetUserByID)
				r.Route("/challenges", func(r chi.Router) {
					r.Get("/", GetChallengeByID)
					r.Route("/:id", func(r chi.Router) {
						r.Get("/", GetChallengeByID)
						r.Get("/pending", GetChallengeByID)
						r.Get("/active", GetChallengeByID)
						r.Get("/completed", GetChallengeByID)
					})
				})
			})
		})

		r.Route("/segments", func(r chi.Router) {
			r.Route("/:id", func(r chi.Router) {
				r.Get("/", GetSegmentByID)
				r.Get("/strava", GetSegmentByIDFromStrava)
			})
			r.Route("/efforts", func(r chi.Router) {
				r.Get("/:id", GetEffortsBySegmentIDFromStrava)
			})
		})

		r.Route("/efforts", func(r chi.Router) {
			r.Get("/:id", GetEffortsBySegmentIDFromStrava)
		})

		r.Route("/challenges", func(r chi.Router) {
			r.Get("/:id", GetChallengeByID)
			r.Post("/accept", GetChallengeByID)
			r.Post("/decline", GetChallengeByID)
			r.Post("/complete", GetChallengeByID)

			r.Route("/pending", func(r chi.Router) {
				r.Get("/:id", GetChallengeByID)
			})

			r.Route("/active", func(r chi.Router) {
				r.Get("/:id", GetChallengeByID)
			})

			r.Route("/completed", func(r chi.Router) {
				r.Get("/:id", GetChallengeByID)
			})
		})

		r.Route("/athletes", func(r chi.Router) {
			r.Route("/:id", func(r chi.Router) {
				r.Get("/", GetAthleteByIDFromStrava)
				r.Get("/friends", GetFriendsByUserIDFromStrava)
				r.Get("/segments", GetSegmentsByUserIDFromStrava)
			})
		})
	})

	mux.Route("/strava", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Get("/", AuthHandler)
		})
	})

	return mux
}
