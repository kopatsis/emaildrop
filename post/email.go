package post

import (
	"emaildrop/database"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendConfirmationEmail(client *sendgrid.Client, entry database.Entry) error {
	from := mail.NewEmail("Me", "donotreply@kopatsis.com")
	to := mail.NewEmail("Me", "j@kopatsis.com")
	subject := "Successful Contact Form Submission Received"
	body := "Here is the data:\n\n" +
		"ID: " + fmt.Sprintf("%d", entry.ID) + "\n" +
		"Timestamp: " + entry.Timestamp.String() + "\n" +
		"RequestID: " + entry.RequestID + "\n" +
		"Name: " + entry.Name + "\n" +
		"Email: " + entry.Email + "\n" +
		"Subject: " + entry.Subject + "\n" +
		"Comment: " + entry.Comment + "\n" +
		"Complete: " + fmt.Sprintf("%v", entry.Complete) + "\n" +
		"IPHash: " + entry.IPHash + "\n" +
		"City: " + entry.City + "\n" +
		"Country: " + entry.Country + "\n" +
		"HasError: " + fmt.Sprintf("%v", entry.HasError) + "\n" +
		"ErrorSource: " + entry.ErrorSource + "\n" +
		"ErrorMessage: " + entry.ErrorMessage + "\n" +
		"Response: " + entry.Response
	message := mail.NewSingleEmail(from, subject, to, body, body)

	_, err := client.Send(message)
	if err != nil {
		return err
	}

	return nil
}

func SendConfirmationToUser(client *sendgrid.Client, name string, email string) error {
	confirmFilePath := filepath.Join(".", "assets", "confirm.html")

	htmlContent, err := os.ReadFile(confirmFilePath)
	if err != nil {
		return err
	}

	from := mail.NewEmail("Demetrios Kopatsis", "donotreply@kopatsis.com")
	to := mail.NewEmail(name, email)
	subject := "I Have Successfully Received Your Contact Form Submission"
	message := mail.NewSingleEmail(from, subject, to, string(htmlContent), string(htmlContent))

	_, err = client.Send(message)
	if err != nil {
		return err
	}

	return nil
}
