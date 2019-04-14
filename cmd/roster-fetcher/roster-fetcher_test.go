package main

import (
	"testing"
	"time"

	"github.com/srabraham/run-director-helper/parkrun"
)

func TestRosterFetch(t *testing.T) {
	ts := parkrun.StartTestServer()
	defer ts.Close()

	prBaseURL = &ts.URL
	now = time.Date(2019, time.March, 1, 0, 0, 0, 0, time.UTC)
	ne := nextEvent()
	if ne.RoleVolunteers[2].Role != "Timekeeper" {
		t.Errorf("Wrong! %v", ne.RoleVolunteers[2])
	}
	if ne.RoleVolunteers[2].Volunteer != "Rod ROMAN" {
		t.Errorf("Wrong! %v", ne.RoleVolunteers[2])
	}
}
