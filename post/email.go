package post

import (
	"emaildrop/database"
	"os"
	"path/filepath"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendConfirmationEmail(client *sendgrid.Client, entry database.Entry) error {
	from := mail.NewEmail("Me", "donotreply@kopatsis.com")
	to := mail.NewEmail("Demetrios Kopatsis", "j@kopatsis.com")
	subject := "Internal: My Server Was Just Contacted"
	htmlBody := "<html><body>" +
		"<p>Hello,</p>" +
		"<p>A new contact form submission was received. Here are the details:</p>" +
		"<p><strong>Request ID:</strong> " + entry.RequestID + "</p>" +
		"<p><strong>Timestamp:</strong> " + entry.Timestamp.Format("2006-01-02 15:04:05") + "</p>" +
		"<p><strong>Name:</strong> " + entry.Name + "</p>" +
		"<p><strong>Email:</strong> " + entry.Email + "</p>" +
		"<p><strong>Subject:</strong> " + entry.Subject + "</p>" +
		"<p><strong>Comment:</strong> " + entry.Comment + "</p>" +
		"<p>If this was not expected, please check for errors in the submission process.</p>" +
		"<p>Best regards,<br>Your Server</p>" +
		"</body></html>"

	message := mail.NewSingleEmail(from, subject, to, "", htmlBody)
	message.AddContent(mail.NewContent("text/html", htmlBody))

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
	subject := "I Got Your Message, and I'll Be in Touch Soon"
	message := mail.NewSingleEmail(from, subject, to, string(htmlContent), string(htmlContent))

	_, err = client.Send(message)
	if err != nil {
		return err
	}

	return nil
}
