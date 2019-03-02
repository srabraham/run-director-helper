package main

import (
	"log"
	"time"

	"github.com/srabraham/run-director-helper/parkrun"
)

func main() {
	roster, err := parkrun.FetchFutureRoster("http://www.parkrun.us/southbouldercreek/futureroster/")
	if err != nil {
		log.Fatal(err)
	}
	nextEvent, err := roster.FirstEventAfter(time.Now())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Next event is\n%v", nextEvent)
}
