package models

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/strava/go.strava"
	"gopkg.in/mgo.v2/bson"
)

// Map struct handles the MongoDB schema for each segments map
type Map struct {
	ID       string `bson:"id" json:"id"`
	Polyline string `bson:"polyline" json:"polyline"`
}

// Segment struct handles the MongoDB schema for a segment
type Segment struct {
	ID                 int64      `bson:"_id" json:"id"`
	Name               string     `bson:"name" json:"name"`
	ActivityType       string     `bson:"activityType" json:"activityType"`
	Distance           float64    `bson:"distance" json:"distance"`
	AverageGrade       float64    `bson:"averageGrade" json:"averageGrade"`
	MaximumGrade       float64    `bson:"maximumGrade" json:"maximumGrade"`
	ElevationHigh      float64    `bson:"elevationHigh" json:"elevationHigh"`
	ElevationLow       float64    `bson:"elevationLow" json:"elevationLow"`
	ClimbCategory      int        `bson:"climbCategory" json:"climbCategory"`
	City               string     `bson:"city" json:"city"`
	State              string     `bson:"state" json:"state"`
	Country            string     `bson:"country" json:"country"`
	TotalElevationGain float64    `bson:"totalElevationGain" json:"totalElevationGain"`
	StartLocation      [2]float64 `bson:"startLocation" json:"startLocation"`
	EndLocation        [2]float64 `bson:"endLocation" json:"endLocation"`
	Map                Map        `bson:"map" json:"map"`
	CreatedAt          time.Time  `bson:"createdAt" json:"createdAt"`
	UpdatedAt          time.Time  `bson:"updatedAt" json:"updatedAt"`
}

// GetSegmentByID gets a single stored segment from MongoDB
func GetSegmentByID(id int64) (*Segment, error) {
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

// SaveSegment stores a cached segment
// this prevents strava api rate limiting
func SaveSegment(s *strava.SegmentDetailed) (*Segment, error) {
	segment := &Segment{
		ID:                 s.Id,
		Name:               s.Name,
		ActivityType:       string(s.ActivityType),
		Distance:           s.Distance,
		AverageGrade:       s.AverageGrade,
		MaximumGrade:       s.MaximumGrade,
		ElevationHigh:      s.ElevationHigh,
		ElevationLow:       s.ElevationLow,
		TotalElevationGain: s.TotalElevationGain,
		ClimbCategory:      int(s.ClimbCategory),
		City:               s.City,
		State:              s.State,
		Country:            s.Country,
		StartLocation:      s.StartLocation,
		EndLocation:        s.EndLocation,
		Map: Map{
			ID:       s.Map.Id,
			Polyline: string(s.Map.Polyline),
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := session.DB(name).C("segments").Insert(&segment); err != nil {
		log.WithField("ID", segment.ID).Errorf("Unable to create segment:\n %v", err)
		return nil, err
	}

	return segment, nil
}

// UpdateSegment stores a cached segment
// this prevents stale data from strava api rate limiting
func (segment Segment) UpdateSegment(s *strava.SegmentDetailed) (*Segment, error) {
	segment.ID = s.Id
	segment.Name = s.Name
	segment.ActivityType = string(s.ActivityType)
	segment.Distance = s.Distance
	segment.AverageGrade = s.AverageGrade
	segment.MaximumGrade = s.MaximumGrade
	segment.ElevationHigh = s.ElevationHigh
	segment.ElevationLow = s.ElevationLow
	segment.TotalElevationGain = s.TotalElevationGain
	segment.ClimbCategory = int(s.ClimbCategory)
	segment.City = s.City
	segment.State = s.State
	segment.Country = s.Country
	segment.StartLocation = s.StartLocation
	segment.EndLocation = s.EndLocation
	segment.Map = Map{
		ID:       s.Map.Id,
		Polyline: string(s.Map.Polyline),
	}
	segment.UpdatedAt = time.Now()

	if err := session.DB(name).C("segments").UpdateId(segment.ID, &segment); err != nil {
		log.WithField("SEGMENT ID", segment.ID).Errorf("Unable to update segment:\n %v", err)
		return nil, err
	}

	return &segment, nil
}
