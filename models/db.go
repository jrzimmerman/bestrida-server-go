package models

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/jrzimmerman/bestrida-server-go/utils"
	"gopkg.in/mgo.v2"
)

// global session to be used in models
var session *mgo.Session
var host = utils.GetEnvString("DB_HOST")
var name = utils.GetEnvString("DB_NAME")
var username = utils.GetEnvString("DB_USER")
var password = utils.GetEnvString("DB_PASSWORD")

// Create new MongoDB session on init
func init() {
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
		logrus.WithError(err).Fatal("Unable to create new session")
	}
}

// Close will close the global MongoDB session
func Close() {
	session.Close()
}
