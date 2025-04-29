package get

import (
	"crypto/sha256"
	"emaildrop/database"
	"emaildrop/middleware"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func hashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

func getSinceParam(c *gin.Context) time.Time {
	sinceStr := c.Query("since")
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return time.Now().Add(-24 * time.Hour)
	}
	if t, err := time.ParseInLocation("2006-01-02T15:04:05", sinceStr, loc); err != nil {
		return time.Now().Add(-24 * time.Hour)
	} else {
		return t
	}
}

func GetResps(tools *middleware.Tools) gin.HandlerFunc {
	return func(c *gin.Context) {
		pass := c.GetHeader("X-Pass")
		if hashPassword(pass) != os.Getenv("PERSONAL_KEY") {
			c.String(http.StatusUnauthorized, "Error: unauthorized")
			return
		}

		since := getSinceParam(c)
		entries, err := database.GetEntriesSince(tools.DB, since)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err.Error()))
			return
		}

		var b strings.Builder
		for _, entry := range entries {
			b.WriteString(fmt.Sprintf(
				`Request ID:    %s
Timestamp:     %s
Name:          %s
Email:         %s
Subject:       %s
Comment:       %s
Complete:      %v
IP Hash:       %s
City:          %s
Country:       %s
Has Error:     %v
Error Source:  %s
Error Message: %s
Response:      %s

`,
				entry.RequestID,
				entry.Timestamp.Format("2006-01-02 15:04:05"),
				entry.Name,
				entry.Email,
				entry.Subject,
				entry.Comment,
				entry.Complete,
				entry.IPHash,
				entry.City,
				entry.Country,
				entry.HasError,
				entry.ErrorSource,
				entry.ErrorMessage,
				entry.Response,
			))
		}

		c.String(http.StatusOK, b.String())
	}
}

func GetResp(tools *middleware.Tools) gin.HandlerFunc {
	return func(c *gin.Context) {
		pass := c.GetHeader("X-Pass")
		if hashPassword(pass) != os.Getenv("PERSONAL_KEY") {
			c.String(http.StatusUnauthorized, "Error: unauthorized")
			return
		}

		reqID := c.Param("id")
		entry, err := database.GetEntry(tools.DB, reqID)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err.Error()))
			return
		}

		text := fmt.Sprintf(
			`Request ID:    %s
Timestamp:     %s
Name:          %s
Email:         %s
Subject:       %s
Comment:       %s
Complete:      %v
IP Hash:       %s
City:          %s
Country:       %s
Has Error:     %v
Error Source:  %s
Error Message: %s
Response:      %s
`,
			entry.RequestID,
			entry.Timestamp.Format("2006-01-02 15:04:05"),
			entry.Name,
			entry.Email,
			entry.Subject,
			entry.Comment,
			entry.Complete,
			entry.IPHash,
			entry.City,
			entry.Country,
			entry.HasError,
			entry.ErrorSource,
			entry.ErrorMessage,
			entry.Response,
		)

		c.String(http.StatusOK, text)
	}
}
