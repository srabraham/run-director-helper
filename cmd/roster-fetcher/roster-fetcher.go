package main

import (
	"fmt"
	"log"

	"github.com/srabraham/gphoto-auto-album-creator/parkrun"
)

func main() {
	result, err := parkrun.FetchFutureRoster("http://www.parkrun.us/southbouldercreek/futureroster/")
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range *result {
		fmt.Println(v)
	}
}
