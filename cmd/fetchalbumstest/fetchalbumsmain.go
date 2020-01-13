package main

import (
	"github.com/srabraham/run-director-helper/cloudfunctions/fetchalbums"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
)

func main() {
	record := httptest.NewRecorder()
	fetchalbums.ListAlbums(record, &http.Request{})
	res, err := ioutil.ReadAll(record.Result().Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("result = %v", string(res))
}
