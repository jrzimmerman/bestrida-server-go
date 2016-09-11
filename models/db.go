package models

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

func getEnvString(env string) string {
	str, ok := os.LookupEnv(env)
	if !ok {
		log.WithField("env", env).Fatal("Missing required environment variable")
	}
	return str
}

// New returns a new session of the MongoDB
func New() (*mgo.Session, error) {
	host := getEnvString("DB_HOST")
	name := getEnvString("DB_NAME")
	username := getEnvString("DB_USER")
	password := getEnvString("DB_PASSWORD")

	// We need this object to establish a session to our MongoDB.
	dbInfo := &mgo.DialInfo{
		Addrs:    []string{host},
		Timeout:  60 * time.Second,
		Database: name,
		Username: username,
		Password: password,
	}

	session, err := mgo.DialWithInfo(dbInfo)
	if err != nil {
		log.WithError(err).Error("Unable to create new session")
	}
	return session, err
}
