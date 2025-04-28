package middleware

import (
	"net/mail"
	"strings"
)

type ContactForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Subject             string `form:"subject"`
	Comment             string `form:"comment"`
	CFTurnstileResponse string `form:"cf-turnstile-response"`
}

func VerifyStruct(form *ContactForm) bool {
	if form == nil {
		return false
	}

	if len(form.Name) == 0 || len(form.Name) > 128 {
		return false
	}
	if len(form.Email) == 0 || len(form.Email) > 320 {
		return false
	}
	form.Email = strings.ToLower(form.Email)

	if len(form.Subject) == 0 || len(form.Subject) > 172 {
		return false
	}
	if len(form.Comment) == 0 || len(form.Comment) > 1024 {
		return false
	}
	if len(form.CFTurnstileResponse) == 0 {
		return false
	}

	return true
}

func VerifyEmail(form *ContactForm, tools *Tools) bool {
	if form == nil || tools == nil {
		return false
	}

	_, err := mail.ParseAddress(form.Email)
	if err != nil {
		return false
	}

	result, err := tools.EmailVerifier.Verify(form.Email)
	if err != nil || result == nil {
		return true
	}

	return result.Syntax.Valid && result.SMTP.Deliverable
}
