package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/srabraham/run-director-helper/googleapis"

	gmail "google.golang.org/api/gmail/v1"
	photoslibrary "google.golang.org/api/photoslibrary/v1"
)

var (
	sharedAlbumName  = flag.String("shared-album-name", "", "Name for shared album to retrieve or create")
	destinationEmail = flag.String("destination-email", "", "Email address to which to send the album link")
)

func createAndShareAlbum(googleClient *http.Client) string {
	photosSvc, err := photoslibrary.New(googleClient)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Creating new album '%s'", *sharedAlbumName)
	resp, err := photosSvc.Albums.Create(
		&photoslibrary.CreateAlbumRequest{
			Album: &photoslibrary.Album{Title: *sharedAlbumName}}).Do()
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

func sendSharingEmail(googleClient *http.Client, shareableURL string) {
	msg := fmt.Sprintf("\r\nCreated a shared album named %s at %s", *sharedAlbumName, shareableURL)
	if err := googleapis.SendEmail(googleClient, "me", "me", *destinationEmail, "Shared album: "+*sharedAlbumName, msg); err != nil {
		log.Fatal(err)
	}
	log.Printf("Successfully sent email")
}

func main() {
	flag.Parse()

	if *sharedAlbumName == "" {
		log.Fatal("Must set a --shared-album-name")
	}
	if *destinationEmail == "" {
		log.Fatal("Must set a --dstEmail")
	}

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
	photosSvc, err := photoslibrary.New(client)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Querying to find any shared albums with the target name")
	resp, err := photosSvc.SharedAlbums.List().Do()
	if err != nil {
		log.Fatal(err)
	}

	containsTargetName := false
	for _, album := range resp.SharedAlbums {
		if album.Title == *sharedAlbumName {
			containsTargetName = true
			break
		}
	}

	if containsTargetName {
		log.Printf("Found album with name '%s'. Aborting...", *sharedAlbumName)
	} else {
		shareableURL := createAndShareAlbum(client)
		sendSharingEmail(client, shareableURL)
	}
}
