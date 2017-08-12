package utils

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// GetEnvString will check for the environment variable, and if found return the string
func GetEnvString(env string) string {
	str, ok := os.LookupEnv(env)
	if !ok {
		log.WithField("ENV", env).Fatal("Missing required environment variable")
	}
	return str
}
