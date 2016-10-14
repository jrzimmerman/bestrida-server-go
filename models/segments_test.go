package models

import "testing"

func TestGetSegmentByID(t *testing.T) {
	id := 2539276

	segment, err := GetSegmentByID(id)
	if err != nil {
		t.Errorf("Unable to retrieve segment by ID:\n %v", err)
	}

	if segment.ID != id {
		t.Errorf("Segment ID %v, is not equal to %v", segment.ID, id)
	}
}
