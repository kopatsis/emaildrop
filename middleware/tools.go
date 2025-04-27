package middleware

import (
	"database/sql"
	"emaildrop/database"
	"log"
	"net/http"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/oschwald/geoip2-golang"
	"github.com/sendgrid/sendgrid-go"
)

type Tools struct {
	SendGrid      *sendgrid.Client
	Client        *http.Client
	Geo           *geoip2.Reader
	EmailVerifier *emailverifier.Verifier
	DB            *sql.DB
}

func (t *Tools) Setup() {
	if db, err := database.SetupDatabase(); err != nil {
		log.Fatal("Database error: " + err.Error())
	} else {
		t.DB = db
	}

}
