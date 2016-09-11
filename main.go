package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/braintree/manners"
	"github.com/jrzimmerman/bestrida-server-go/routes"
)

func main() {
	port := getEnvString("PORT")
	log.WithField("PORT", port).Info("Listening for http traffic")
	log.Fatal(manners.ListenAndServe(":"+port, routes.API()))
}

func getEnvString(env string) string {
	str, ok := os.LookupEnv(env)
	if !ok {
		log.WithField("env", env).Fatal("Missing required environment variable")
	}
	return str
}
