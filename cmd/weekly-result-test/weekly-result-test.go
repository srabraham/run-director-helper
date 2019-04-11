package main

import (
	"flag"
	"log"
	"time"

	"github.com/srabraham/run-director-helper/parkrun"
)

var (
	latestResultsURL = flag.String("latest-results-url", "http://www.parkrun.us/southbouldercreek/results/latestresults/", "URL for a parkrun latest results page")
)

func importantNumber(n int) bool {
	return n == 10 || n == 25 || n%50 == 0
}

func main() {
	ne, err := parkrun.NextEventNumber(*latestResultsURL)
	if err != nil {
		log.Fatal(err)
	}
	eventNum := ne - 1
	idToRecentRunCount := make(map[int64]int)
	idToRunner := make(map[int64]parkrun.Runner)
	for i := 0; i < 5; i++ {
		eventNum--
		time.Sleep(time.Second)
		log.Printf("Getting results for event %d", eventNum)
		er, err := parkrun.GetRunners("http://www.parkrun.us/southbouldercreek/results/weeklyresults", int32(eventNum))
		if err != nil {
			log.Fatal(err)
		}
		for _, r := range er.Runners {
			idToRecentRunCount[r.AthleteID]++
			idToRunner[r.AthleteID] = r
		}
	}
	for id, rCount := range idToRecentRunCount {
		if importantNumber(rCount + 1) {
			log.Printf("Milestone: %s (%d) has %d runs", idToRunner[id].Name, idToRunner[id].AthleteID, idToRunner[id].TotalRuns)
		} else if importantNumber(rCount+2) || importantNumber(rCount+3) {
			log.Printf("Near milestone: %s (%d) has %d runs", idToRunner[id].Name, idToRunner[id].AthleteID, idToRunner[id].TotalRuns)
		}
	}
}
