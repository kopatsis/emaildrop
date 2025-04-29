package post

import (
	"bytes"
	"emaildrop/database"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func SendConfirmationSlack(client *http.Client, entry database.Entry) error {

	webhookURL := os.Getenv("SLACK_WEBHOOK")

	text := fmt.Sprintf(
		"*New Contact Form Submission:*\nRequest ID: %s\nTimestamp: %s\nName: %s\nEmail: %s\nSubject: %s\nComment: %s\nComplete: %v",
		entry.RequestID,
		entry.Timestamp.String(),
		entry.Name,
		entry.Email,
		entry.Subject,
		entry.Comment,
		entry.Complete,
	)

	payload := map[string]string{
		"text": text,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
