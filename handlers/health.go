package handlers

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// GetHealthCheck returns response for health check
func GetHealthCheck(w http.ResponseWriter, r *http.Request) {
	res := New(w)
	body := map[string]interface{}{"status": "healthy"}
	log.Info("healthy")
	res.Render(http.StatusOK, body)
}
