package main

import (
	"emaildrop/middleware"
	"emaildrop/post"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	t := middleware.Tools{}
	t.Setup()

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RateLimitMiddleware(t.Limiter))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello\nMy site is here: https://kopatsis.com")
	})

	r.POST("/contact", post.PostContactForm(&t))

	r.Run(":8080")
}
