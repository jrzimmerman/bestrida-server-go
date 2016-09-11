package models

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// Map struct handles the MongoDB schema for each segments map
type Map struct {
	ID            string `json:"id"`
	Polyline      string `json:"polyline"`
	ResourceState int    `json:"resourceState"`
}

// Segment struct handles the MongoDB schema for a segment
type Segment struct {
	ID                 int       `json:"_id"`
	ResourceState      int       `json:"resourceState"`
	Name               string    `json:"name"`
	ActivityType       string    `json:"activityType"`
	Distance           float64   `json:"distance"`
	AverageGrade       float64   `json:"averageGrade"`
	MaximumGrade       float64   `json:"maximumGrade"`
	ClimbCategory      float64   `json:"climbCategory"`
	City               string    `json:"city"`
	State              string    `json:"string"`
	Country            string    `json:"country"`
	TotalElevationGain float64   `json:"totalElevationGain"`
	EndLatLng          []float64 `json:"endLatLng"`
	StartLatLng        []float64 `json:"startLatLng"`
	Map                Map       `json:"map"`
}

// GetSegmentByID gets a single stored segment from MongoDB
func GetSegmentByID(id int) (*Segment, error) {
	var seg Segment

	s, err := New()
	if err != nil {
		log.WithError(err).Error("Unable to create new MongoDB Session")
		return nil, err
	}

	defer s.Close()

	err = s.DB("heroku_zgxbr4j2").C("segments").Find(bson.M{"_id": id}).One(&seg)
	if err != nil {
		log.WithField("ID", id).Error("Unable to find segment with id")
		return nil, err
	}

	return &seg, nil
}
