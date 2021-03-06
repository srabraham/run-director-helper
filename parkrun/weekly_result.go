package parkrun

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

const (
	weeklyResultsURL = "/results/weeklyresults"
	headerPos        = "Pos"
	headerParkrunner = "parkrunner"
	headerTime       = "Time"
	headerAgeCat     = "Age Cat"
	headerAgeGrade   = "Age Grade"
	headerGender     = ""
	headerGenderPos  = "Gender Pos"
	headerClub       = "Club"
	headerNote       = "Note"
	headerTotalRuns  = "Total Runs"
)

var (
	expectedHeader       = []string{headerPos, headerParkrunner, headerTime, headerAgeCat, headerAgeGrade, headerGender, headerGenderPos, headerClub, headerNote, headerTotalRuns}
	athleteNumberMatcher = regexp.MustCompile("athleteNumber=([0-9]+)")
)

type Runner struct {
	Name      string
	AthleteID int64
	TotalRuns int
}

type EventRunners struct {
	eventNum int32
	Runners  []Runner
}

func GetRunners(prBaseURL string, eventNum int32) (*EventRunners, error) {
	resp, err := http.Get(fmt.Sprintf("%s%s/?runSeqNumber=%d", prBaseURL, weeklyResultsURL, eventNum))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to get weekly result:\n%v", resp)
	}
	defer resp.Body.Close()
	return getRunners(resp.Body, eventNum)
}

func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			log.Printf("%s != %s", a[i], b[i])
			return false
		}
	}
	return true
}

func getRunners(html io.Reader, eventNum int32) (*EventRunners, error) {
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return nil, err
	}
	var resultsErr error
	er := EventRunners{
		Runners: make([]Runner, 0),
	}
	doc.Find("#results").Each(func(i int, s *goquery.Selection) {
		headers := make([]string, 0)

		s.Find("thead tr th").Each(func(i int, s *goquery.Selection) {
			headers = append(headers, s.Text())
		})

		if !slicesEqual(headers, expectedHeader) {
			resultsErr = fmt.Errorf("Headers != expected headers: [%v] != [%v]", headers, expectedHeader)
			return
		}
		headerByIndex := make(map[int]string)
		for i, h := range headers {
			headerByIndex[i] = h
		}

		s.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
			r := Runner{}
			s.Find("td").Each(func(i int, s *goquery.Selection) {
				switch headerByIndex[i] {
				case headerParkrunner:
					r.Name = s.Text()
					if s.Text() != "Unknown" {
						html, err := s.Html()
						if err != nil {
							resultsErr = err
							return
						}
						matches := athleteNumberMatcher.FindStringSubmatch(html)
						if len(matches) < 2 {
							resultsErr = fmt.Errorf("Failed to find athleteNumber in %s", html)
							return
						}
						r.AthleteID, err = strconv.ParseInt(matches[1], 10, 64)
						if err != nil {
							resultsErr = err
							return
						}
					}
				case headerTotalRuns:
					if len(s.Text()) > 0 {
						tr, err := strconv.Atoi(s.Text())
						if err != nil {
							resultsErr = err
							return
						}
						r.TotalRuns = tr
					}
				}
			})
			if r.Name != "Unknown" {
				er.Runners = append(er.Runners, r)
			}
		})
	})
	if resultsErr != nil {
		return nil, resultsErr
	}
	er.eventNum = eventNum
	return &er, nil
}
