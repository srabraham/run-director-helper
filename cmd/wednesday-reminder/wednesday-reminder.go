package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/srabraham/run-director-helper/googleapis"
	"github.com/srabraham/run-director-helper/parkrun"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/photoslibrary/v1"
)

var (
	futureRosterURL = flag.String("future-roster-url", "http://www.parkrun.us/southbouldercreek/futureroster/", "URL for a parkrun future roster page")
)

func main() {
	flag.Parse()

	if err := googleapis.AddScope(gmail.GmailSendScope,
		photoslibrary.PhotoslibraryAppendonlyScope,
		photoslibrary.PhotoslibraryReadonlyScope,
		photoslibrary.PhotoslibrarySharingScope); err != nil {
		log.Fatal(err)
	}
	_, err := googleapis.GetClient()
	if err != nil {
		log.Fatal(err)
	}

	roster, err := parkrun.FetchFutureRoster(*futureRosterURL)
	if err != nil {
		log.Fatal(err)
	}
	nextEvent, err := roster.FirstEventAfter(time.Now())
	if err != nil {
		log.Fatal(err)
	}
	nextRd := nextEvent.VolunteersForRole("Run Director")
	dateStr := nextEvent.Date.Format("2006-01-02")
	nextRdStr := strings.Trim(fmt.Sprint(nextRd), "[]")
	subject := "Update for run on " + dateStr
	message := fmt.Sprintf(
		"Hi run director %s,\n"+
			"Here is the roster for the upcoming run:\n"+
			"%v", nextRdStr, nextEvent)
	log.Printf("Would use subject\n%s", subject)
	log.Printf("Would send message\n%s", message)
}
