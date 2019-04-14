package main

import (
	"testing"
	"time"

	"github.com/srabraham/run-director-helper/parkrun"
)

func TestMilestoneFinder(t *testing.T) {
	ts := parkrun.StartTestServer()
	defer ts.Close()

	prBaseURL = &ts.URL
	zeroTime := time.Duration(0)
	timeBetweenGets = &zeroTime

	// TODO: finish the test!
	// _ = getUpcomingMilestones()
}
