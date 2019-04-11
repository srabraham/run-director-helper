package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/srabraham/google-oauth-helper/googleauth"
	"github.com/srabraham/run-director-helper/googleapis"
	"github.com/srabraham/run-director-helper/parkrun"
	"google.golang.org/api/gmail/v1"
)

var (
	prBaseURL           = flag.String("pr-base-url", "http://www.parkrun.us/southbouldercreek", "Base URL for parkrun event")
	necessaryVolunteers = flag.String("necessary-volunteers", "Run Director,Equipment Storage and Delivery,Timekeeper,Barcode Scanning,Finish Tokens,Marshal,Marshal", "Comma-separated list of required volunteer positions, possibly including duplicates")
	destinationEmail    = flag.String("destination-email", "", "Email address to which to send the reminder")
)

func missingVolunteers(nextEvent parkrun.EventDetails) string {
	necessaryVolunteersCounts := make(map[string]int)
	for _, v := range strings.Split(*necessaryVolunteers, ",") {
		necessaryVolunteersCounts[v]++
	}
	for _, rv := range nextEvent.RoleVolunteers {
		// Empty string implies the role is vacant
		if len(rv.Volunteer) == 0 {
			continue
		}
		if necessaryVolunteersCounts[rv.Role] > 0 {
			necessaryVolunteersCounts[rv.Role]--
		}
	}
	missingVolunteersStrs := make([]string, 0)
	for k, v := range necessaryVolunteersCounts {
		if v > 0 {
			missingVolunteersStrs = append(missingVolunteersStrs, fmt.Sprintf("%d %s", v, k))
		}
	}
	missingVolunteersMsg := strings.Join(missingVolunteersStrs, ", ")
	return missingVolunteersMsg
}

func main() {
	flag.Parse()
	if *destinationEmail == "" {
		log.Fatal("Must set a destination-email")
	}

	if err := googleauth.AddScope(gmail.GmailSendScope); err != nil {
		log.Fatal(err)
	}
	if err := googleauth.SetTokenFileName("gmailsend-tok"); err != nil {
		log.Fatal(err)
	}
	googleClient, err := googleauth.GetClient()
	if err != nil {
		log.Fatal(err)
	}

	roster, err := parkrun.FetchFutureRoster(*prBaseURL)
	if err != nil {
		log.Fatal(err)
	}
	nextEvent, err := roster.FirstEventAfter(time.Now())
	if err != nil {
		log.Fatal(err)
	}

	missingVolunteersMsg := missingVolunteers(nextEvent)

	nextRd := nextEvent.VolunteersForRole("Run Director")
	dateStr := nextEvent.Date.Format("2006-01-02")
	nextRdStr := strings.Trim(fmt.Sprint(nextRd), "[]")

	subject := ""
	if len(missingVolunteersMsg) > 0 {
		subject = fmt.Sprintf("[Automated] Volunteers needed for %v event", dateStr)
	} else {
		subject = fmt.Sprintf("[Automated] No additional volunteers needed for %v event", dateStr)
	}

	message := fmt.Sprintf("Hi run director %s,\n\n", nextRdStr)
	if len(missingVolunteersMsg) > 0 {
		message += fmt.Sprintf("We still need volunteers for these roles: %s\n\n", missingVolunteersMsg)
	} else {
		message += "We have all the required roles filled for the next run!\n\n"
	}
	message += fmt.Sprintf("Here is the roster as of now:\n%v\n\n", nextEvent)
	message += fmt.Sprintf("You can see the roster on the web at %s", parkrun.FutureRosterURL(*prBaseURL))

	log.Printf("Email subject\n%s", subject)
	log.Printf("Email message\n%s", message)

	if err = googleapis.SendEmail(googleClient, "me", "me", *destinationEmail, subject, message); err != nil {
		log.Fatal(err)
	}
}
