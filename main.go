package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/braintree/manners"
	"github.com/jrzimmerman/bestrida-server-go/handlers"
	"github.com/jrzimmerman/bestrida-server-go/models"
)

func main() {
	port := getEnvString("PORT")
	log.WithField("PORT", port).Info("Listening for http traffic")

	// close DB connection
	defer models.Close()

	log.Fatal(manners.ListenAndServe(":"+port, handlers.API()))
}

func getEnvString(env string) string {
	str, ok := os.LookupEnv(env)
	if !ok {
		log.WithField("ENV", env).Fatal("Missing required environment variable")
	}
	return str
}
