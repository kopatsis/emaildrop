package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

func IPData(tools *Tools, c *gin.Context) (string, string, string) {
	var city string
	var country string

	ipStr := c.ClientIP()

	if ipStr == "" || ipStr == "::1" {
		ipStr = c.Request.Header.Get("X-Forwarded-For")
	}

	hashed := ""
	if ipStr != "" {
		if commaIndex := strings.Index(ipStr, ","); commaIndex != -1 {
			ipStr = ipStr[:commaIndex]
		}

		ip := net.ParseIP(ipStr)
		if ip != nil {
			record, err := tools.Geo.City(ip)
			if err == nil && record != nil {
				city = record.City.Names["en"]
				country = record.Country.Names["en"]
			}
		}
		hashed = hex.EncodeToString(sha256.New().Sum([]byte(ipStr)))
	}

	return city, country, hashed

}
