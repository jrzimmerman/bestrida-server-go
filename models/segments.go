package models

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// Map struct handles the MongoDB schema for each segments map
type Map struct {
	ID            string `bson:"id" json:"id"`
	Polyline      string `bson:"polyline" json:"polyline"`
	ResourceState int    `bson:"resource_state" json:"resource_state"`
}

// Segment struct handles the MongoDB schema for a segment
type Segment struct {
	ID                 int       `bson:"_id" json:"_id"`
	ResourceState      int       `bson:"resourceState"json:"resourceState"`
	Name               string    `bson:"name" json:"name"`
	ActivityType       string    `bson:"activityType" json:"activityType"`
	Distance           float64   `bson:"distance" json:"distance"`
	AverageGrade       float64   `bson:"averageGrade" json:"averageGrade"`
	MaximumGrade       float64   `bson:"maximumGrade" json:"maximumGrade"`
	ClimbCategory      float64   `bson:"climbCategory" json:"climbCategory"`
	City               string    `bson:"city" json:"city"`
	State              string    `bson:"state" json:"state"`
	Country            string    `bson:"country" json:"country"`
	TotalElevationGain float64   `bson:"totalElevationGain" json:"totalElevationGain"`
	EndLatLng          []float64 `bson:"endLatLng" json:"endLatLng"`
	StartLatLng        []float64 `bson:"startLatLng" json:"startLatLng"`
	Map                Map       `bson:"map" json:"map"`
}

// GetSegmentByID gets a single stored segment from MongoDB
func GetSegmentByID(id int) (*Segment, error) {
	var s Segment

	if err := session.DB(name).C("segments").Find(bson.M{"_id": id}).One(&s); err != nil {
		log.WithField("ID", id).Error("Unable to find segment with id")
		return nil, err
	}

	log.WithFields(map[string]interface{}{
		"NAME": s.Name,
		"ID":   s.ID,
	}).Info("Segment returned from DB")
	return &s, nil
}
