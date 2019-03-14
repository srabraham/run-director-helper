package main

import (
	"testing"

	"github.com/srabraham/run-director-helper/parkrun"
)

func TestNecessaryVolunteersFullRoster(t *testing.T) {
	*necessaryVolunteers = "Cheerleader,Cheerleader,Lead dog"
	str := missingVolunteers(parkrun.EventDetails{RoleVolunteers: []parkrun.RoleVolunteer{
		{Role: "Cheerleader", Volunteer: "Jim"},
		{Role: "Cheerleader", Volunteer: "Bob"},
		{Role: "Lead dog", Volunteer: "Wicket"},
	}})
	if len(str) > 0 {
		t.Errorf("Expected necessaryVolunteers string to be empty, was: %s", str)
	}
}

func TestNecessaryVolunteersEmptyStringVolunteer(t *testing.T) {
	*necessaryVolunteers = "Cheerleader,Cheerleader,Cheerleader"
	str := missingVolunteers(parkrun.EventDetails{RoleVolunteers: []parkrun.RoleVolunteer{
		{Role: "Cheerleader", Volunteer: ""},
	}})
	expect := "3 Cheerleader"
	if str != expect {
		t.Errorf("Expected necessaryVolunteers string to be %s, was: %s", expect, str)
	}
}

func TestNecessaryVolunteersNoEntryForVolunteer(t *testing.T) {
	*necessaryVolunteers = "Cheerleader"
	str := missingVolunteers(parkrun.EventDetails{RoleVolunteers: []parkrun.RoleVolunteer{}})
	expect := "1 Cheerleader"
	if str != expect {
		t.Errorf("Expected necessaryVolunteers string to be %s, was: %s", expect, str)
	}
}
