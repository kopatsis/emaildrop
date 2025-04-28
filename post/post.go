package post

import (
	"emaildrop/database"
	"emaildrop/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func PostContactForm(tools *middleware.Tools) gin.HandlerFunc {
	return func(c *gin.Context) {

		entry := database.Entry{
			RequestID: "RQ-" + uuid.NewString(),
			Timestamp: time.Now(),
		}

		var form middleware.ContactForm
		if err := c.ShouldBind(&form); err != nil || !middleware.VerifyStruct(&form, tools) {
			entry.HasError = true
			entry.ErrorSource = "Binding Form"
			if err != nil {
				entry.ErrorMessage = err.Error()
			} else {
				entry.ErrorMessage = "Failed on entry struct verification"
			}

		}
	}
}
