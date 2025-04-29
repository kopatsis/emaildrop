package main

import (
	"emaildrop/get"
	"emaildrop/middleware"
	"emaildrop/post"
	"net/http"

	// "os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	r := gin.Default()

	// isProduction := os.Getenv("ENV") == "production"
	// if isProduction {
	// 	trustedProxies := []string{
	// 		"kopatsis.com",
	// 		"192.168.1.100",
	// 	}
	// 	r.SetTrustedProxies(trustedProxies)
	// } else {
	// 	r.SetTrustedProxies([]string{"127.0.0.1", "localhost"})
	// }

	t := middleware.Tools{}
	t.Setup()

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RateLimitMiddleware(t.Limiter))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello\nMy site is here: https://kopatsis.com")
	})

	r.POST("/contact", post.PostContactForm(&t))

	r.GET("/data", get.GetResps(&t))
	r.GET("/data/:id", get.GetResp(&t))

	r.Run(":8080")
}
