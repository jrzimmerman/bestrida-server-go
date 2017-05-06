package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/jrzimmerman/bestrida-server-go/handlers"
	"github.com/jrzimmerman/bestrida-server-go/models"
	"github.com/jrzimmerman/bestrida-server-go/utils"
)

func main() {
	port := utils.GetEnvString("PORT")
	log.WithField("PORT", port).Info("Listening for http traffic")

	// close DB connection
	defer models.Close()

	log.Fatal(http.ListenAndServe(":"+port, handlers.API()))
}
