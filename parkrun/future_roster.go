package parkrun

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	futureRosterURL = "/futureroster/"
	// DateFormat is the format used on the future roster web pages.
	DateFormat = "2 January 2006"
)

var (
	eventLocation = flag.String("location", "America/Denver", "Time zone in which the parkrun occurs")
	eventTime     = flag.Duration("event-time", 9*time.Hour, "The duration after midnight local time when the event is held")
)

// RoleVolunteer is a pair of a role and a volunteer.
// The Volunteer may be empty if no volunteer is yet assigned.
type RoleVolunteer struct {
	Role      string
	Volunteer string
}

// EventDetails gives a run date and a list of the volunteers for that run.
type EventDetails struct {
	Date           time.Time
	RoleVolunteers []RoleVolunteer
}

// FutureRoster is a list of events scraped from a future roster web page.
type FutureRoster struct {
	SortedEvents []EventDetails
}

// FirstEventAfter finds the first event in the roster after the provided time.
func (fr FutureRoster) FirstEventAfter(t time.Time) (EventDetails, error) {
	for _, v := range fr.SortedEvents {
		if v.Date.After(t) {
			return v, nil
		}
	}
	return EventDetails{}, errors.New("Found no events on future roster after " + t.String())
}

// VolunteersForRole returns the volunteer(s) for the provided role name.
func (details EventDetails) VolunteersForRole(role string) []string {
	volunteers := make([]string, 0)
	for _, rv := range details.RoleVolunteers {
		if rv.Role == role && rv.Volunteer != "" {
			volunteers = append(volunteers, rv.Volunteer)
		}
	}
	return volunteers
}

func (details EventDetails) String() string {
	roleNames := make([]string, 0)
	sort.Strings(roleNames)
	str := fmt.Sprintf("%s [\n", details.Date.Format("2006-01-02"))
	for _, rv := range details.RoleVolunteers {
		str += fmt.Sprintf("  %s: %s\n", rv.Role, rv.Volunteer)
	}
	str += "]"
	return str
}

func FutureRosterURL(basePrURL string) string {
	return basePrURL + futureRosterURL
}

// FetchFutureRoster gets the volunteer rosters from the provided URL.
func FetchFutureRoster(basePrURL string) (FutureRoster, error) {
	resp, e := http.Get(basePrURL + futureRosterURL)
	if e != nil {
		return FutureRoster{}, e
	}
	defer resp.Body.Close()
	return fetchFutureRoster(resp.Body)
}

func fetchFutureRoster(html io.Reader) (FutureRoster, error) {
	loc, err := time.LoadLocation(*eventLocation)
	if err != nil {
		return FutureRoster{}, err
	}
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return FutureRoster{}, err
	}

	roster := make([]EventDetails, 0)
	doc.Find("#rosterTable").Each(func(i int, s *goquery.Selection) {
		headers := make([]string, 0)
		rows := make([][]string, 0)

		s.Find("thead tr th").Each(func(i int, s *goquery.Selection) {
			headers = append(headers, s.Text())
		})

		s.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
			row := make([]string, 0)
			role := s.Find("th a").Text()
			row = append(row, role)
			s.Find("td").Each(func(i int, s *goquery.Selection) {
				volunteer := strings.Trim(s.Text(), "\n ")
				row = append(row, volunteer)
			})
			rows = append(rows, row)
		})

		for i := 1; i < len(headers); i++ {
			t, err := time.ParseInLocation(DateFormat, headers[i], loc)
			if err != nil {
				log.Fatal(err)
			}
			t = t.Add(*eventTime)
			rv := make([]RoleVolunteer, 0)
			for j := 0; j < len(rows); j++ {
				volunteer := rows[j][i]
				role := rows[j][0]
				rv = append(rv, RoleVolunteer{Role: role, Volunteer: volunteer})
			}
			roster = append(roster, EventDetails{Date: t, RoleVolunteers: rv})
		}

	})
	if len(roster) == 0 {
		return FutureRoster{}, errors.New("couldn't find a roster")
	}
	return FutureRoster{SortedEvents: roster}, nil
}
