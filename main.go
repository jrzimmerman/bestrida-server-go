package main

import (
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/jrzimmerman/bestrida-server-go/routes"
)

func main() {
	port := getEnvString("PORT")
	logrus.WithField("port", port).Info("Listening for http traffic")
	logrus.Fatal(http.ListenAndServe(":"+port, routes.API()))
}

func getEnvString(env string) string {
	str, ok := os.LookupEnv(env)
	if !ok {
		logrus.WithField("env", env).Fatal("Missing required environment variable")
	}
	return str
}
