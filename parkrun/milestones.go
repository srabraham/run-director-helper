package parkrun

import (
	"fmt"
	"log"
	"time"
)

func GetUpcomingMilestones(prBaseURL string, maxPastEvents int, timeBetweenGets time.Duration) []Runner {
	ne, err := NextEventNumber(prBaseURL)
	if err != nil {
		log.Fatal(err)
	}
	eventNum := ne - 1
	idToRecentRunCount := make(map[int64]int)
	idToRunner := make(map[int64]Runner)
	for i := 0; i < maxPastEvents; i++ {
		eventNum--
		if eventNum <= 0 {
			log.Print("Breaking due to hitting event 1")
			break
		}
		time.Sleep(timeBetweenGets)
		log.Printf("Getting results for event %d", eventNum)
		er, err := GetRunners(prBaseURL, int32(eventNum))
		if err != nil {
			log.Fatal(err)
		}
		for _, r := range er.Runners {
			idToRecentRunCount[r.AthleteID]++
			idToRunner[r.AthleteID] = r
		}
	}
	milestones := make([]Runner, 0)
	for id, rCount := range idToRecentRunCount {
		txt := fmt.Sprintf("%d runs for %s (%d)", idToRunner[id].TotalRuns, idToRunner[id].Name, idToRunner[id].AthleteID)
		if importantNumber(rCount + 1) {
			fmt.Printf("Milestone: %s\n", txt)
			milestones = append(milestones, idToRunner[id])
		} else if importantNumber(rCount+2) || importantNumber(rCount+3) {
			fmt.Printf("Near milestone: %s\n", txt)
		} else {
			fmt.Printf("Not near milestone: %s\n", txt)
		}
	}
	return milestones
}

func importantNumber(n int) bool {
	return n == 10 || n == 25 || n%50 == 0
}
