package fetchalbums

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"cloud.google.com/go/firestore"
)

type Album struct {
	url  string
	num  int64
	name string
}

func ListAlbums(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, err := firestore.NewClient(ctx, "sbcparkrun")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	albumYears, err := client.Collection("albumyears").Documents(ctx).GetAll()
	if err != nil {
		log.Fatal(err)
	}

	allAlbums := make([]Album, 0)
	for _, albumYear := range albumYears {
		for _, album := range albumYear.Data() {
			albumDecode := album.(map[string]interface{})
			allAlbums = append(allAlbums,
				Album{
					name: albumDecode["name"].(string),
					num:  albumDecode["num"].(int64),
					url:  albumDecode["url"].(string),
				})
		}
	}
	sort.Slice(allAlbums, func(i, j int) bool {
		return allAlbums[i].num > allAlbums[j].num
	})
	fmt.Fprintf(w, "<!doctype html>\n")
	fmt.Fprintf(w, "<html>\n")
	fmt.Fprintf(w, "<body>\n")
	fmt.Fprintf(w, "  <ul style=\"font-family:Lato, sans-serif;\">\n")
	for _, album := range allAlbums {
		fmt.Fprintf(w, "    <li><a href=%s target=\"_blank\">%s</a></li>\n", album.url, album.name)
	}
	fmt.Fprintf(w, "  </ul>\n")
	fmt.Fprintf(w, "</body>\n")
	fmt.Fprintf(w, "</html>\n")
}
