package middleware

import (
	"database/sql"
	"emaildrop/database"
	"log"
	"net/http"
	"os"
	"time"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/oschwald/geoip2-golang"
	"github.com/sendgrid/sendgrid-go"
)

type Tools struct {
	SendGrid      *sendgrid.Client
	Client        *http.Client
	Limiter       *limiter.Limiter
	Geo           *geoip2.Reader
	EmailVerifier *emailverifier.Verifier
	DB            *sql.DB
}

func (t *Tools) Setup() {
	if t == nil {
		return
	}

	limiter := tollbooth.NewLimiter(5, &limiter.ExpirableOptions{
		DefaultExpirationTTL: time.Minute,
	})
	t.Limiter = limiter

	if db, err := database.SetupDatabase(); err != nil {
		log.Fatal("Database error: " + err.Error())
	} else {
		t.DB = db
	}

	if sgAPIKey := os.Getenv("SENDGRID_API_KEY"); sgAPIKey != "" {
		t.SendGrid = sendgrid.NewSendClient(sgAPIKey)
	} else {
		log.Fatal("Missing SENDGRID_API_KEY")
	}

	t.Client = &http.Client{
		Timeout: 30 * time.Second,
	}

	if geo, err := geoip2.Open("assets/GeoLite2-City.mmdb"); err != nil {
		log.Fatal("GeoIP2 error: " + err.Error())
	} else {
		t.Geo = geo
	}

	t.EmailVerifier = emailverifier.NewVerifier()
}
