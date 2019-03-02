package main

import (
	"fmt"
	"log"
	"time"

	"github.com/srabraham/run-director-helper/parkrun"
)

func main() {
	result, err := parkrun.FetchFutureRoster("http://www.parkrun.us/southbouldercreek/futureroster/")
	if err != nil {
		log.Fatal(err)
	}
	var nextEvent parkrun.EventDetails
	for _, v := range result {
		if v.Date.After(time.Now()) {
			nextEvent = v
			break
		}
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
