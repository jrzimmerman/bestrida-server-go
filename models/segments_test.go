package models

import "testing"

func TestGetSegmentByIDSuccess(t *testing.T) {
	id := 2539276

	segment, err := GetSegmentByID(id)
	if err != nil {
		t.Errorf("Unable to retrieve segment by ID:\n %v", err)
	}

	if segment.ID != id {
		t.Errorf("Segment ID %v, is not equal to %v", segment.ID, id)
	}
}

func TestGetSegmentByIDFailure(t *testing.T) {
	id := 0

	_, err := GetSegmentByID(id)
	if err.Error() != "not found" {
		t.Errorf("Unable to throw error for ID:\n %v", err)
	}
}
