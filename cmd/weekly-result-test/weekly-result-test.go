package main

import (
	"flag"
	"log"

	"github.com/srabraham/run-director-helper/parkrun"
)

var (
	prBaseURL = flag.String("pr-base-url", "http://www.parkrun.us/southbouldercreek", "Base parkrun URL")
)

func main() {
	ms := parkrun.GetUpcomingMilestones(*prBaseURL, 5)
	log.Printf("Milestones = %v", ms)
}
