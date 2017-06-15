package handlers

import (
	"fmt"
	"net/http"
)

// LandingHandler serves the static landing content for the landing page
func LandingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "landing")
}
