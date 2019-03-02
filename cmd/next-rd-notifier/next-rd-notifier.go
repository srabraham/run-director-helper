package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/srabraham/run-director-helper/parkrun"
)

var (
	futureRosterURL = flag.String("future-roster-url", "http://www.parkrun.us/southbouldercreek/futureroster/", "URL for a parkrun future roster page")
)

func main() {
	roster, err := parkrun.FetchFutureRoster(*futureRosterURL)
	if err != nil {
		log.Fatal(err)
	}
	nextEvent, err := roster.FirstEventAfter(time.Now())
	if err != nil {
		log.Fatal(err)
	}
	nextRd := nextEvent.VolunteersForRole("Run Director")
	log.Printf("Next event is\n%v", nextEvent)
	var message string
	if len(nextRd) == 0 {
		message = fmt.Sprintf(
			"WARNING: No run director for parkrun on %v",
			nextEvent.Date.Format("2006-01-02"))
	} else {
		message = fmt.Sprintf(
			"%v will be run director for %v",
			nextRd,
			nextEvent.Date.Format("2006-01-02"))
	}
	log.Printf("Message to send: %s", message)
}
