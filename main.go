package main

import (
	"emaildrop/get"
	"emaildrop/middleware"
	"emaildrop/post"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	r := gin.Default()

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
