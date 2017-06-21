package utils

import (
	"os"

	"github.com/Sirupsen/logrus"
)

// GetEnvString will check for the environment variable, and if found return the string
func GetEnvString(env string) string {
	str, ok := os.LookupEnv(env)
	if !ok {
		logrus.WithField("ENV", env).Fatal("Missing required environment variable")
	}
	return str
}
