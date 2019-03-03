package googleapis

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/api/gmail/v1"
)

func SendEmail(authedClient *http.Client, senderUserID, from, to, subject, message string) error {
	gmailSvc, err := gmail.New(authedClient)
	if err != nil {
		return err
	}
	messageToEncode := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"\r\n%s",
		from, to, subject, message)
	log.Printf("Will send email:\n%s", messageToEncode)
	rawEmail := base64.StdEncoding.EncodeToString([]byte(messageToEncode))
	_, err = gmail.NewUsersMessagesService(gmailSvc).Send(senderUserID,
		&gmail.Message{Raw: rawEmail}).Do()
	if err != nil {
		return err
	}
	log.Printf("Successfully sent email")
	return nil
}
