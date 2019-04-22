package parkrun

import (
	"fmt"
	"log"
	"sort"
	"time"
)

func GetUpcomingMilestones(prBaseURL string, maxPastEvents int, timeBetweenGets time.Duration) []Runner {
	le, err := LastEventNumber(prBaseURL)
	if err != nil {
		log.Fatal(err)
	}
	eventNum := le
	idToRunner := make(map[int64]Runner)
	for i := 0; i < maxPastEvents; i++ {
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
			idToRunner[r.AthleteID] = r
		}
		eventNum--
	}
	sortedRunners := make([]Runner, 0, len(idToRunner))
	for _, r := range idToRunner {
		sortedRunners = append(sortedRunners, r)
	}
	sort.Slice(sortedRunners, func(i, j int) bool {
		if sortedRunners[i].TotalRuns != sortedRunners[j].TotalRuns {
			return sortedRunners[i].TotalRuns < sortedRunners[j].TotalRuns
		}
		return sortedRunners[i].AthleteID < sortedRunners[j].AthleteID
	})
	milestones := make([]Runner, 0)
	for _, r := range sortedRunners {
		txt := fmt.Sprintf("%d runs for %s (%d)", r.TotalRuns, r.Name, r.AthleteID)
		if importantNumber(r.TotalRuns + 1) {
			fmt.Printf("One run until milestone: %s\n", txt)
			milestones = append(milestones, r)
		} else if importantNumber(r.TotalRuns+2) {
			fmt.Printf("Two runs until milestone: %s\n", txt)
		} else {
			fmt.Printf("Not near milestone: %s\n", txt)
		}
	}
	return milestones
}

func importantNumber(n int) bool {
	return n == 10 || n%25 == 0
}
