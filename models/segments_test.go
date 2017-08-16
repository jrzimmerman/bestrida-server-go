package models

import (
	"testing"

	"github.com/jrzimmerman/bestrida-server-go/utils"
	log "github.com/sirupsen/logrus"
	strava "github.com/strava/go.strava"
)

func TestGetSegmentByIDSuccess(t *testing.T) {
	id := int64(2539276)

	segment, err := GetSegmentByID(id)
	if err != nil {
		t.Errorf("Unable to retrieve segment by ID:\n %v", err)
	}

	if segment.ID != id {
		t.Errorf("Segment ID %v, is not equal to %v", segment.ID, id)
	}
}

func TestGetSegmentByIDFailure(t *testing.T) {
	id := int64(0)

	_, err := GetSegmentByID(id)
	if err.Error() != "not found" {
		t.Errorf("Unable to throw error for ID:\n %v", err)
	}
}

var accessToken = utils.GetEnvString("STRAVA_ACCESS_TOKEN")

func TestSaveSegment(t *testing.T) {
	var numID int64 = 2539276

	if err := RemoveSegment(numID); err != nil {
		t.Error("Unable to remove segment")
		return
	}
	t.Logf("segment %d removed", numID)

	// use our access token to grab generic segment info
	client := strava.NewClient(accessToken)
	segment, err := strava.NewSegmentsService(client).Get(numID).Do()
	if err != nil {
		log.Error("Unable to retrieve segment info")
		return
	}
	t.Logf("segment %d returned from Strava", segment.Id)

	stored, err := SaveSegment(segment)
	if err != nil {
		t.Errorf("Unable to store segment from Strava:\n %v", err)
		return
	}

	if stored.ID != segment.Id {
		t.Errorf("Stored segment ID %v, is not equal to %v from Strava", stored.ID, segment.Id)
		return
	}
	t.Logf("segment %d successfully stored", stored.ID)

	dup, err := SaveSegment(segment)
	if err == nil {
		t.Errorf("Should throw error when trying to store duplicate segment: %v", dup.ID)
	}
}

func TestUpdateSegment(t *testing.T) {
	var numID int64 = 2539276

	// use our access token to grab generic segment info
	client := strava.NewClient(accessToken)
	stravaSegment, err := strava.NewSegmentsService(client).Get(numID).Do()
	if err != nil {
		log.Error("Unable to retrieve segment info")
		return
	}
	t.Logf("segment %d returned from Strava", stravaSegment.Id)

	segment, err := GetSegmentByID(numID)
	if err != nil {
		t.Errorf("Unable to retrieve segment by ID:\n %v", err)
	}

	if segment.ID != numID {
		t.Errorf("Segment ID %v, is not equal to %v", segment.ID, numID)
	}

	updated, err := segment.UpdateSegment(stravaSegment)
	if err != nil {
		t.Errorf("Unable to update stored segment:\n %v", err)
		return
	}

	if updated.ID != stravaSegment.Id {
		t.Errorf("Updated segment ID %v, is not equal to %v from Strava", updated.ID, stravaSegment.Id)
		return
	}
	t.Logf("segment %d successfully updated", updated.ID)
}
