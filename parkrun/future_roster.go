package parkrun

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	dateFormat = "2 January 2006"
)

type RoleVolunteer struct {
	Role      string
	Volunteer string
}

type EventDetails struct {
	Date           time.Time
	RoleVolunteers []RoleVolunteer
}

func (details EventDetails) String() string {
	roleNames := make([]string, 0)
	sort.Strings(roleNames)
	str := fmt.Sprintf("%s [\n", details.Date.Format(dateFormat))
	for _, rv := range details.RoleVolunteers {
		str += fmt.Sprintf("  %s: %s\n", rv.Role, rv.Volunteer)
	}
	str += "]"
	return str
}

func FetchFutureRoster(url string) (*[]EventDetails, error) {

	doc, err := goquery.NewDocument(url)
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
				volunteer := s.Text()
				row = append(row, volunteer)
			})
			rows = append(rows, row)
		})

		for i := 1; i < len(headers); i++ {
			t, err := time.Parse(dateFormat, headers[i])
			if err != nil {
				log.Fatal(err)
			}
			rv := make([]RoleVolunteer, 0)
			for j := 0; j < len(rows); j++ {
				volunteer := rows[j][i]
				role := rows[j][0]
				rv = append(rv, RoleVolunteer{Role: role, Volunteer: volunteer})
			}
			roster = append(roster, EventDetails{Date: t, RoleVolunteers: rv})
		}

	})
	return &roster, nil
}
