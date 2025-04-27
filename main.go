package main

import (
	"emaildrop/middleware"
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

	r.POST("/contact", func(c *gin.Context) {
		var form middleware.ContactForm
		if err := c.ShouldBind(&form); err == nil {
			// You can process the form here later
			c.JSON(http.StatusOK, gin.H{
				"message": "Form received",
				"data":    form,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})

	r.Run(":8080")
}
