package main

import (
	"flag"
	"log"
	"time"

	"github.com/srabraham/run-director-helper/parkrun"
)

var (
	prBaseURL = flag.String("pr-base-url", "http://www.parkrun.us/southbouldercreek", "Base URL for parkrun event")
	now       = time.Now()
)

func nextEvent() parkrun.EventDetails {
	roster, err := parkrun.FetchFutureRoster(*prBaseURL)
	if err != nil {
		log.Fatal(err)
	}
	nextEvent, err := roster.FirstEventAfter(now)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Next event is\n%v", nextEvent)
	return nextEvent
}

func main() {
	nextEvent()
}
