package middleware

import (
	"crypto/sha256"
	"emaildrop/database"
	"encoding/hex"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

func IPData(tools *Tools, c *gin.Context, entry *database.Entry) {
	if tools == nil || c == nil || entry == nil {
		return
	}

	ipStr := c.ClientIP()

	if ipStr == "" || ipStr == "::1" {
		ipStr = c.Request.Header.Get("X-Forwarded-For")
	}

	if ipStr != "" {
		if commaIndex := strings.Index(ipStr, ","); commaIndex != -1 {
			ipStr = ipStr[:commaIndex]
		}

		ip := net.ParseIP(ipStr)
		if ip != nil {
			record, err := tools.Geo.City(ip)
			if err == nil && record != nil {
				entry.City = record.City.Names["en"]
				entry.Country = record.Country.Names["en"]
			}
		}
		entry.IPHash = hex.EncodeToString(sha256.New().Sum([]byte(ipStr)))
	}
}
