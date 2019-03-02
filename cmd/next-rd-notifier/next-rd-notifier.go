package main

import (
	"fmt"
	"log"
	"time"

	"github.com/srabraham/run-director-helper/parkrun"
)

func main() {
	result, err := parkrun.FetchFutureRoster("http://www.parkrun.us/southbouldercreek/futureroster/")
	if err != nil {
		log.Fatal(err)
	}
	time.Date(result[0].Date)
	for _, v := range result {
		fmt.Println(v)
	}
}
