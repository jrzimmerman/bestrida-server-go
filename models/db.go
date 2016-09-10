package models

import (
	"os"

	"github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

var session *mgo.Session

// New returns a new session of the MongoDB
func New() *mgo.Session {
	connection, ok := os.LookupEnv("DB_CONN")
	if !ok {
		logrus.WithField("DB_CONN", connection).Fatal("Database connection not passed as environment variable")
	}
	session, err := mgo.Dial(connection)
	if err != nil {
		logrus.WithError(err).Fatal("Unable to connnect to database")
	}
	return session
}

// Close will close the MongoDB session
func Close() {
	session.Close()
}
