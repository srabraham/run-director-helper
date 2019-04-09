package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/srabraham/google-oauth-helper/googleauth"
	"github.com/srabraham/run-director-helper/googleapis"
	"github.com/srabraham/run-director-helper/parkrun"

	docs "google.golang.org/api/docs/v1"
	gmail "google.golang.org/api/gmail/v1"
	photoslibrary "google.golang.org/api/photoslibrary/v1"
)

var (
	destinationEmail = flag.String("destination-email", "", "Email address to which to send the album link")
	futureRosterURL  = flag.String("future-roster-url", "http://www.parkrun.us/southbouldercreek/futureroster/", "URL for a parkrun future roster page")
	latestResultsURL = flag.String("latest-results-url", "http://www.parkrun.us/southbouldercreek/results/latestresults/", "URL for a parkrun latest results page")
	albumDocID       = flag.String("album-doc-id", "1fCvOX4sUiKOrXvuRE9pd0K40qFO1eCaP0wxoRF1I2YY", "ID of a Google Doc ID that will contain the album links")
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
	yesterday := time.Now().Add(-24 * time.Hour)
	// Use yesterday. This allows the program to be run up to 24 hours an event start,
	// while still referring to that event rather than next week's event.
	nextEvent, err := fr.FirstEventAfter(yesterday)
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

func updateDoc(googleClient *http.Client, albumName string, shareableURL string) {
	docsSvc, err := docs.New(googleClient)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	newText := albumName + "\n"
	// Put the new link on the first line.
	insertionIndex := int64(1)
	resp, err := docsSvc.Documents.BatchUpdate(
		*albumDocID,
		&docs.BatchUpdateDocumentRequest{
			Requests: []*docs.Request{
				{InsertText: &docs.InsertTextRequest{
					Location: &docs.Location{
						Index: insertionIndex,
					},
					Text: newText,
				}},
				{UpdateTextStyle: &docs.UpdateTextStyleRequest{
					Fields: "link",
					Range: &docs.Range{
						StartIndex: insertionIndex,
						EndIndex:   insertionIndex + int64(len(newText)),
					},
					TextStyle: &docs.TextStyle{
						Link: &docs.Link{
							Url: shareableURL,
						},
					},
				}},
			},
		}).Do()
	log.Printf("Resp = %v", resp)
	log.Println("Updated the Google Doc with the new album URL")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	if *destinationEmail == "" {
		log.Fatal("Must set a --destination-email")
	}

	albumName := getAlbumName()

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
		updateDoc(client, albumName, shareableURL)
	}
}
