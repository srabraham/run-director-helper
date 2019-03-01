package main

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

type roleVolunteer struct {
	role      string
	volunteer string
}

type runDetails struct {
	date           time.Time
	roleVolunteers []roleVolunteer
}

func (details runDetails) String() string {
	roleNames := make([]string, 0)
	sort.Strings(roleNames)
	str := fmt.Sprintf("%s [\n", details.date.Format(dateFormat))
	for _, rv := range details.roleVolunteers {
		str += fmt.Sprintf("  %s: %s\n", rv.role, rv.volunteer)
	}
	str += "]"
	return str
}

func scrape(url string) (*[]runDetails, error) {

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	roster := make([]runDetails, 0)
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
			rv := make([]roleVolunteer, 0)
			for j := 1; j < len(rows); j++ {
				volunteer := rows[j][i]
				role := rows[j][0]
				rv = append(rv, roleVolunteer{role: role, volunteer: volunteer})
			}
			roster = append(roster, runDetails{date: t, roleVolunteers: rv})
		}

	})
	return &roster, nil
}

func main() {
	result, err := scrape("http://www.parkrun.us/southbouldercreek/futureroster/")
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range *result {
		fmt.Println(v)
	}
}
