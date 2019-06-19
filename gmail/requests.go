// Package gmail handels authentication and retrieval of emails from gmail
package gmail

import (
	"fmt"
	"log"

	"google.golang.org/api/gmail/v1"
)

// GetEmail get a single email
func GetEmail(id string) *gmail.Message {
	user := "me"
	srv := getService()
	msg, err := srv.Users.Messages.Get(user, id).Format("full").Do()
	if err != nil {
		log.Printf("Error: gmail.go/GetMessage/msg returned %v\n", err)
	}
	return msg
}

// GetEmails get meany emails with one request
// this is a temporary solution until I can
// check for a batch get option in the gmail api
func GetEmails(IDs []string) {
	for _, id := range IDs {
		GetEmail(id)
	}
}

// GetEmailBody finds the payload that is html and returns it in it's original urlbase64 encoding
func GetEmailBody(msg *gmail.Message) string {
	for i := 0; i < len(msg.Payload.Parts); i++ {
		if msg.Payload.Parts[i].MimeType == "text/html" {
			return msg.Payload.Parts[i].Body.Data
		}

		for j := 0; j < len(msg.Payload.Parts[i].Parts); j++ {
			if msg.Payload.Parts[i].Parts[j].MimeType == "text/html" {
				return msg.Payload.Parts[i].Parts[j].Body.Data
			}
		}
	}
	return ""
}

// ListEmails lists all emails matching a label and query
func ListEmails(labelID, query string) []*gmail.Message {
	user := "me"
	srv := getService()

	r, err := srv.Users.Messages.List(user).LabelIds(labelID).Q(query).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve emails: %v", err)
	}

	return r.Messages
}

// ListLables all lables and their ids
func ListLables() {
	user := "me"
	srv := getService()

	r, err := srv.Users.Labels.List(user).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}
	if len(r.Labels) == 0 {
		fmt.Println("No labels found.")
		return
	}
	fmt.Println("Labels:")
	for _, l := range r.Labels {
		fmt.Printf("- %s\tID: %s\n", l.Name, l.Id)
	}
}
