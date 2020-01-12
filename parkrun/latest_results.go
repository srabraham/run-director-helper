package parkrun

import (
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

const (
	latestURL = "/results/latestresults/"
)

var (
	parkrunNumberMatcher = regexp.MustCompile(`^#([0-9]+)$`)
)

func LastEventNumber(prBaseURL string) (int64, error) {
	resp, err := http.Get(prBaseURL + latestURL)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()
	return lastEventNumber(resp.Body)
}

func lastEventNumber(html io.Reader) (int64, error) {
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return -1, err
	}
	var lastRunNumber int64
	doc.Find("div").Each(func(i int, s *goquery.Selection) {
		if s.HasClass("Results-header") {
			s.Find("span").Each(func(i1 int, s1 *goquery.Selection) {
				matches := parkrunNumberMatcher.FindStringSubmatch(s1.Text())
				if len(matches) == 2 {
					lastRunNumber, err = strconv.ParseInt(matches[1], 10, 64)
					if err != nil {
						log.Printf("Failed to find parkrun number in %s", s.Text())
						return
					}
				}
			})
		}
	})
	return lastRunNumber, nil
}
