package handlers

import (
	"github.com/jrzimmerman/bestrida-server-go/utils"
	"github.com/pressly/chi"
	strava "github.com/strava/go.strava"
)

var authenticator *strava.OAuthAuthenticator
var clientID = utils.GetEnvString("STRAVA_CLIENT_ID")
var clientSecret = utils.GetEnvString("STRAVA_CLIENT_SECRET")
var accessToken = utils.GetEnvString("STRAVA_ACCESS_TOKEN")
var port = utils.GetEnvString("PORT")

// API initializes all endpoints
func API() (mux *chi.Mux) {
	mux = chi.NewRouter()
	mux.Use(CORS)

	mux.Route("/api", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Get("/:id", GetUserByID)
		})
		r.Route("/segments", func(r chi.Router) {
			r.Get("/:id", GetSegmentByID)
		})
		r.Route("/challenges", func(r chi.Router) {
			r.Get("/:id", GetChallengeByID)
		})

		r.Route("/strava", func(r chi.Router) {
			r.Route("/auth", func(r chi.Router) {
				r.Get("/", AuthHandler)
			})
			r.Route("/athletes", func(r chi.Router) {
				r.Route("/:id", func(r chi.Router) {
					r.Get("/", GetAthleteByIDFromStrava)

					r.Route("/friends", func(r chi.Router) {
						r.Get("/", GetFriendsByUserIDFromStrava)
					})

					r.Route("/segments", func(r chi.Router) {
						r.Get("/", GetSegmentsByUserIDFromStrava)
						r.Route("/:segmentID", func(r chi.Router) {
							r.Get("/", GetSegmentByIDFromStrava)
							r.Route("/efforts", func(r chi.Router) {
								r.Get("/", GetEffortsBySegmentIDFromStrava)
							})
						})
					})
				})
			})

			r.Route("/segments", func(r chi.Router) {
				r.Get("/:id", GetSegmentByIDFromStrava)
			})
		})
	})

	return mux
}
