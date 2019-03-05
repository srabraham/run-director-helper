package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/srabraham/run-director-helper/parkrun"

	"github.com/srabraham/run-director-helper/googleapis"

	gmail "google.golang.org/api/gmail/v1"
	photoslibrary "google.golang.org/api/photoslibrary/v1"
)

var (
	destinationEmail = flag.String("destination-email", "", "Email address to which to send the album link")
	futureRosterURL  = flag.String("future-roster-url", "http://www.parkrun.us/southbouldercreek/futureroster/", "URL for a parkrun future roster page")
	latestResultsURL = flag.String("latest-results-url", "http://www.parkrun.us/southbouldercreek/results/latestresults/", "URL for a parkrun latest resutls page")
)

func getAlbumIDIfExists(googleClient *http.Client, albumName string) string {
	log.Print("Querying to find any shared albums with the target name")
	photosSvc, err := photoslibrary.New(googleClient)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := photosSvc.SharedAlbums.List().Do()
	if err != nil {
		log.Fatal(err)
	}

	for _, album := range resp.SharedAlbums {
		if album.Title == albumName {
			return album.Id
		}
	}
	return ""
}

func createAndShareAlbum(googleClient *http.Client, albumName string) string {
	photosSvc, err := photoslibrary.New(googleClient)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Creating new album '%s'", albumName)
	resp, err := photosSvc.Albums.Create(
		&photoslibrary.CreateAlbumRequest{
			Album: &photoslibrary.Album{Title: albumName}}).Do()
	if err != nil {
		log.Fatal(err)
	}
	newAlbumID := resp.Id
	log.Printf("Created new album. ID: %s", newAlbumID)
	shareResponse, err := photosSvc.Albums.Share(
		newAlbumID,
		&photoslibrary.ShareAlbumRequest{
			SharedAlbumOptions: &photoslibrary.SharedAlbumOptions{
				IsCollaborative: true,
				IsCommentable:   true,
			}}).Do()
	if err != nil {
		log.Fatal(err)
	}
	shareableURL := shareResponse.ShareInfo.ShareableUrl
	log.Printf("Successfully shared album: %s", shareableURL)
	return shareableURL
}

func sendSharingEmail(googleClient *http.Client, albumName string, shareableURL string) {
	msg := fmt.Sprintf("The next event's shared album, %s, has shareable URL %s", albumName, shareableURL)
	if err := googleapis.SendEmail(googleClient, "me", "me", *destinationEmail, "[Automated] Shared album: "+albumName, msg); err != nil {
		log.Fatal(err)
	}
	log.Printf("Successfully sent email")
}

func getAlbumName() string {
	fr, err := parkrun.FetchFutureRoster(*futureRosterURL)
	if err != nil {
		log.Fatal(err)
	}
	nextEvent, err := fr.FirstEventAfter(time.Now())
	if err != nil {
		log.Fatal(err)
	}
	nextEventDateStr := nextEvent.Date.Format("2006-01-02")
	log.Printf("Next event is on %s", nextEventDateStr)
	lastEventNumber, err := parkrun.NextEventNumber(*latestResultsURL)
	if err != nil {
		log.Fatal(err)
	}
	nextEventNumber := lastEventNumber + 1
	return fmt.Sprintf("SBC parkrun #%d (%s)", nextEventNumber, nextEventDateStr)
}

func main() {
	flag.Parse()

	if *destinationEmail == "" {
		log.Fatal("Must set a --destination-email")
	}

	albumName := getAlbumName()

	if err := googleapis.AddScope(gmail.GmailSendScope,
		photoslibrary.PhotoslibraryAppendonlyScope,
		photoslibrary.PhotoslibraryReadonlyScope,
		photoslibrary.PhotoslibrarySharingScope); err != nil {
		log.Fatal(err)
	}
	if err := googleapis.SetTokenFileName("sharedalbummaker-tok"); err != nil {
		log.Fatal(err)
	}
	client, err := googleapis.GetClient()
	if err != nil {
		log.Fatal(err)
	}

	if getAlbumIDIfExists(client, albumName) != "" {
		log.Printf("Found album with name '%s'. Aborting...", albumName)
	} else {
		shareableURL := createAndShareAlbum(client, albumName)
		sendSharingEmail(client, albumName, shareableURL)
	}
}
