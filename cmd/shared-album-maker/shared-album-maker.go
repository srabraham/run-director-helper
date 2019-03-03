package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"

	"github.com/srabraham/run-director-helper/googleapis"

	gmail "google.golang.org/api/gmail/v1"
	photoslibrary "google.golang.org/api/photoslibrary/v1"
)

var (
	sharedAlbumName = flag.String("shared-album-name", "", "Name for shared album to retrieve or create")
	dstEmail        = flag.String("dst-email", "", "Destination email address")
)

func createAlbumAndSendEmail(photosSvc *photoslibrary.Service, gmailSvc *gmail.Service) {
	log.Printf("New to create new album '%s'", *sharedAlbumName)
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

	rawEmail := base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf(
			"From: 'me'\r\n"+
				"To: %s\r\n"+
				"Subject: Shared album: %s\r\n"+
				"\r\nCreated a shared album named %s at %s",
			*dstEmail, *sharedAlbumName, *sharedAlbumName, shareableURL)))

	_, err = gmail.NewUsersMessagesService(gmailSvc).Send("me",
		&gmail.Message{Raw: rawEmail}).Do()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Successfully sent email")
}

func main() {
	flag.Parse()

	if *sharedAlbumName == "" {
		log.Fatal("Must set a --shared-album-name")
	}
	if *dstEmail == "" {
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
	gmailSvc, err := gmail.New(client)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Querying for shared albums with the target name")
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
		createAlbumAndSendEmail(photosSvc, gmailSvc)
	}
}
