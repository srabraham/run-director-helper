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
	DateFormat = "2 January 2006"
)

var (
	// The duration after midnight local time when the event is held.
	parkrunTime = 9 * time.Hour
	// The timezone in which the event is held.
	location = flag.String("location", "America/Denver", "Time zone in which the parkrun occurs")
)

// RoleVolunteer is a pair of a role and a volunteer.
type RoleVolunteer struct {
	Role      string
	Volunteer string
}

// EventDetails gives a run date and a list of the volunteers for that run.
type EventDetails struct {
	Date           time.Time
	RoleVolunteers []RoleVolunteer
}

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
	str := fmt.Sprintf("%s [\n", details.Date.Format(time.RFC3339))
	for _, rv := range details.RoleVolunteers {
		str += fmt.Sprintf("  %s: %s\n", rv.Role, rv.Volunteer)
	}
	str += "]"
	return str
}

// FetchFutureRoster gets the volunteer rosters from the provided URL.
func FetchFutureRoster(url string) ([]EventDetails, error) {
	resp, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	defer resp.Body.Close()
	return fetchFutureRoster(resp.Body)
}

func fetchFutureRoster(html io.Reader) ([]EventDetails, error) {
	loc, err := time.LoadLocation(*location)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return nil, err
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
			t = t.Add(parkrunTime)
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
		return nil, errors.New("couldn't find a roster")
	}
	return roster, nil
}
