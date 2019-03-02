package main

import (
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
	log.Printf("Next event is\n%v", nextEvent)
}