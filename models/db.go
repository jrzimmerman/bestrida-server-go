package models

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

// global session to be used in models
var session *mgo.Session

// Create new MongoDB session on init
func init() {
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
	var err error
	session, err = mgo.DialWithInfo(dbInfo)
	if err != nil {
		log.WithError(err).Fatal("Unable to create new session")
	}
}

func getEnvString(env string) string {
	str, ok := os.LookupEnv(env)
	if !ok {
		log.WithField("env", env).Fatal("Missing required environment variable")
	}
	return str
}

// Close will close the global MongoDB session
func Close() {
	session.Close()
}
