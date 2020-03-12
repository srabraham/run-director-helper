package main

import (
	"flag"
	"fmt"
	"github.com/srabraham/run-director-helper/albumdb"
	"log"
	"net/http"
	"time"

	"github.com/srabraham/google-oauth-helper/googleauth"
	"github.com/srabraham/run-director-helper/googleapis"
	"github.com/srabraham/run-director-helper/googleapis/generatedapivendor/photoslibrary/v1"
	"github.com/srabraham/run-director-helper/parkrun"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/gmail/v1"
)

var (
	destinationEmail    = flag.String("destination-email", "", "Email address to which to send the album link")
	prBaseURL           = flag.String("pr-base-url", "http://www.parkrun.us/southbouldercreek", "Base URL for parkrun event")
	dbProjectId         = flag.String("db-project-id", "sbcparkrun", "Name of Firestore GCP project.")
	dbCollectionName    = flag.String("db-collection-name", "albumyears", "Name of Firestore collection in which to save albums")
	eventNumberOverride = flag.Int64("event-number-override", 0, "Override for event number")
	eventDateOverride   = flag.String("event-date-override", "", "Override for event date")
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

func getNextEventNumber() int64 {
	if *eventNumberOverride != 0 {
		return *eventNumberOverride
	}
	lastEventNumber, err := parkrun.LastEventNumber(*prBaseURL)
	if err != nil {
		log.Fatal(err)
	}
	return lastEventNumber + 1
}

func getNextEventDate() string {
	if *eventDateOverride != "" {
		return *eventDateOverride
	}
	fr, err := parkrun.FetchFutureRoster(*prBaseURL)
	if err != nil {
		log.Fatal(err)
	}
	yesterday := time.Now().Add(-24 * time.Hour)
	// Use yesterday. This allows the program to be run up to 24 hours an event start,
	// while still referring to that event rather than next week's event.
	nextEvent, err := fr.FirstEventAfter(yesterday)
	if err != nil {
		log.Fatal(err)
	}
	return nextEvent.Date.Format("2006-01-02")
}

func getAlbumName(nextEventNumber int64, nextEventDate string) string {
	return fmt.Sprintf("SBC parkrun #%d (%s)", nextEventNumber, nextEventDate)
}

func main() {
	flag.Parse()

	if *destinationEmail == "" {
		log.Fatal("Must set a --destination-email")
	}

	// Will fail the program if firebase connection can't be made.
	albumdb.TestDbConnection(*dbProjectId)

	nextEventNumber := getNextEventNumber()
	nextEventDate := getNextEventDate()
	albumName := getAlbumName(nextEventNumber, nextEventDate)

	if err := googleauth.AddScope(gmail.GmailSendScope,
		photoslibrary.PhotoslibraryAppendonlyScope,
		photoslibrary.PhotoslibraryReadonlyScope,
		photoslibrary.PhotoslibrarySharingScope,
		docs.DocumentsScope); err != nil {
		log.Fatal(err)
	}
	if err := googleauth.SetTokenFileName("sharedalbummaker-tok"); err != nil {
		log.Fatal(err)
	}
	client, err := googleauth.GetClient()
	if err != nil {
		log.Fatal(err)
	}

	if getAlbumIDIfExists(client, albumName) != "" {
		log.Printf("Found album with name '%s'. Aborting...", albumName)
	} else {
		shareableURL := createAndShareAlbum(client, albumName)
		sendSharingEmail(client, albumName, shareableURL)
		albumdb.AddAlbumToDb(nextEventDate, albumName, shareableURL, *dbProjectId, *dbCollectionName, nextEventNumber)
	}
}
