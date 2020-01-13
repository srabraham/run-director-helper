package main

import "github.com/srabraham/run-director-helper/albumdb"

func main() {
	albumdb.AddAlbumToDb(
		"2021-01-01",
		"my up  updated album 123",
		"http://sean.run",
		"sbcparkrun",
		"albumyears-test",
		1245)
}