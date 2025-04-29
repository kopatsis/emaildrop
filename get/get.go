package get

import (
	"crypto/sha256"
	"emaildrop/database"
	"emaildrop/middleware"
	"encoding/hex"
	"net/http"
	"os"
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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no auth key provided"})
			return
		}

		since := getSinceParam(c)
		if entries, err := database.GetEntriesSince(tools.DB, since); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error_message": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"results": entries})
		}
	}
}

func GetResp(tools *middleware.Tools) gin.HandlerFunc {
	return func(c *gin.Context) {
		pass := c.GetHeader("X-Pass")
		if hashPassword(pass) != os.Getenv("PERSONAL_KEY") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no auth key provided"})
			return
		}

		reqID := c.Param("id")
		if entry, err := database.GetEntry(tools.DB, reqID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error_message": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"results": entry})
		}
	}
}
