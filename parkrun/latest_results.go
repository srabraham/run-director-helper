package parkrun

import (
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

var (
	parkrunNumberMatcher = regexp.MustCompile("(?s)^.+\\sparkrun\\s#[[:space:]]*([0-9]+).*$")
)

func NextEventNumber(url string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()
	return nextEventNumber(resp.Body)
}

func nextEventNumber(html io.Reader) (int, error) {
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return -1, err
	}
	lastRunNumber := 0
	doc.Find("h2").Each(func(i int, s *goquery.Selection) {
		matches := parkrunNumberMatcher.FindStringSubmatch(s.Text())
		if len(matches) < 2 {
			log.Printf("Failed to find parkrun number in %s", s.Text())
			lastRunNumber = -1
			return
		}
		lastRunNumber, err = strconv.Atoi(matches[1])
		if err != nil {
			log.Printf("Failed to convert %s to an int", matches[1])
			lastRunNumber = -1
			return
		}
	})
	return lastRunNumber, nil
}
