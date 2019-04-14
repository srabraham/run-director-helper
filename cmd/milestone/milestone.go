package main

import (
	"flag"
	"log"
	"time"

	"github.com/srabraham/run-director-helper/parkrun"
)

var (
	prBaseURL       = flag.String("pr-base-url", "http://www.parkrun.us/southbouldercreek", "Base parkrun URL")
	numPastEvents   = flag.Int("num-past-events", 5, "Number of past events to query")
	timeBetweenGets = flag.Duration("time-between-gets", time.Second*5, "Time between GETs of weekly results")
)

func getUpcomingMilestones() []parkrun.Runner {
	ms := parkrun.GetUpcomingMilestones(*prBaseURL, *numPastEvents, *timeBetweenGets)
	log.Printf("Milestones = %v", ms)
	return ms
}

func main() {
	getUpcomingMilestones()
}
