package cloudflare

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"time"
)

func VerifyTurnstile(responseToken string, client *http.Client) (bool, error) {
	secretKey := os.Getenv("CF_SECRET_KEY")
	if secretKey == "" {
		return false, errors.New("missing CF_SECRET_KEY environment variable")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	form := url.Values{}
	form.Set("secret", secretKey)
	form.Set("response", responseToken)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://challenges.cloudflare.com/turnstile/v0/siteverify", nil)
	if err != nil {
		return false, err
	}

	req.Form = form

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result struct {
		Success bool `json:"success"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return false, err
	}

	return result.Success, nil
}
