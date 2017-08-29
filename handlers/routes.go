package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/jrzimmerman/bestrida-server-go/utils"
	log "github.com/sirupsen/logrus"
	"github.com/strava/go.strava"
)

var authenticator = &strava.OAuthAuthenticator{
	CallbackURL:            "http://www.bestridaapp.com/strava/auth/callback",
	RequestClientGenerator: nil,
}
var clientID = utils.GetEnvString("STRAVA_CLIENT_ID")
var clientSecret = utils.GetEnvString("STRAVA_CLIENT_SECRET")
var accessToken = utils.GetEnvString("STRAVA_ACCESS_TOKEN")
var port = utils.GetEnvString("PORT")

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

// API initializes all endpoints
func API() (mux *chi.Mux) {
	path, err := authenticator.CallbackPath()
	if err != nil {
		// possibly that the callback url set above is invalid
		log.Errorf("unable to set strava callback path: \n %v", err)
	}
	clientIDInt, err := strconv.Atoi(clientID)
	if err != nil {
		log.Errorf("unable to convert strava client id to int: \n %v", err)
	}
	strava.ClientId = clientIDInt
	strava.ClientSecret = clientSecret

	mux = chi.NewRouter()
	mux.Use(CORS)

	mux.HandleFunc(path, authenticator.HandlerFunc(oAuthSuccess, oAuthFailure))

	workDir, _ := os.Getwd()
	publicDir := filepath.Join(workDir, "public")
	FileServer(mux, "/", http.Dir(publicDir))

	mux.Route("/api", func(r chi.Router) {
		r.Get("/health", GetHealthCheck)
		r.Route("/users", func(r chi.Router) {
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", GetUserByID)
				r.Get("/friends", GetFriendsByUserID)

				r.Route("/segments", func(r chi.Router) {
					r.Get("/", GetSegmentsByUserID)
					r.Route("/{segmentID}", func(r chi.Router) {
						r.Get("/", GetSegmentByIDWithUserID)
						r.Get("/strava", GetSegmentByIDFromStravaWithUserID)
						r.Route("/efforts", func(r chi.Router) {
							r.Get("/", GetEffortsBySegmentIDFromStravaWithUserID)
						})
					})
				})

				r.Route("/challenges", func(r chi.Router) {
					r.Get("/", GetAllChallengesByUserID)
					r.Get("/pending", GetPendingChallengesByUserID)
					r.Get("/active", GetActiveChallengesByUserID)
					r.Get("/completed", GetCompletedChallengesByUserID)
				})

			})
		})

		r.Route("/segments", func(r chi.Router) {
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", GetSegmentByID)
				r.Get("/strava", GetSegmentByIDFromStrava)
			})
		})

		r.Route("/challenges", func(r chi.Router) {
			r.Get("/{id}", GetChallengeByID)
			r.Put("/accept", AcceptChallengeByID)
			r.Put("/decline", DeclineChallengeByID)
			r.Put("/complete", CompleteChallengeByID)
			r.Post("/create", CreateChallenge)
		})

		r.Route("/athletes", func(r chi.Router) {
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", GetAthleteByIDFromStrava)
				r.Get("/friends", GetFriendsByUserIDFromStrava)
				r.Get("/segments", GetSegmentsByUserIDFromStrava)
			})
		})
	})

	mux.Route("/strava", func(r chi.Router) {
		r.Route("/update", func(r chi.Router) {
			r.Get("/users", UpdateAllUsersFromStrava)
		})
		r.Route("/auth", func(r chi.Router) {
			r.Get("/", AuthHandler)
			r.Get("/callback", AuthHandler)
		})
	})

	return mux
}
