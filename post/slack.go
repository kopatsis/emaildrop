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
		"*Request ID:* %s\n*Link:* %s\n*Timestamp:* %s\n*Name:* %s\n*Email:* %s\n*Subject:* %s\n*Comment:* %s\n*Complete:* %v",
		entry.RequestID,
		"https://emaildrop.kopatsis.com/data/"+entry.RequestID,
		entry.Timestamp.Format("2006-01-02 15:04:05"),
		entry.Name,
		entry.Email,
		entry.Subject,
		entry.Comment,
		entry.Complete,
	)

	if entry.IPHash != "" {
		text += fmt.Sprintf("\n*IP Hash:* %s", entry.IPHash)
	}
	if entry.City != "" {
		text += fmt.Sprintf("\n*City:* %s", entry.City)
	}
	if entry.Country != "" {
		text += fmt.Sprintf("\n*Country:* %s", entry.Country)
	}
	if entry.HasError {
		text += fmt.Sprintf("\n*Has Error:* %v", entry.HasError)
	}
	if entry.ErrorSource != "" {
		text += fmt.Sprintf("\n*Error Source:* %s", entry.ErrorSource)
	}
	if entry.ErrorMessage != "" {
		text += fmt.Sprintf("\n*Error Message:* %s", entry.ErrorMessage)
	}
	if entry.Response != "" {
		text += fmt.Sprintf("\n*Response:* %s", entry.Response)
	}

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
