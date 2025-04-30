package post

import (
	"emaildrop/database"
	"emaildrop/middleware"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func applyError(entry *database.Entry, tools *middleware.Tools, c *gin.Context, source string, err error) {
	if entry == nil || c == nil {
		return
	}

	fmt.Println("Error at: " + source + " - " + err.Error())

	entry.HasError = true
	entry.ErrorSource = source
	if err != nil {
		entry.ErrorMessage = err.Error()
	} else {
		entry.ErrorMessage = "no specifc error given"
	}
	entry.Response = "Sorry, I'm unable to process your form submission request at the moment. Please email me directly, thanks :)."

	if err := database.InsertEntry(tools.DB, *entry); err != nil {
		fmt.Printf("Unable to save entry: %+v\n", entry)
		fmt.Printf("\n\n%v\n", err)
	}

}

func PostContactForm(tools *middleware.Tools) gin.HandlerFunc {
	return func(c *gin.Context) {

		entry := database.Entry{
			RequestID: "RQ-" + uuid.NewString(),
			Timestamp: time.Now(),
		}
		middleware.IPData(tools, c, &entry)

		var form middleware.ContactForm
		if err := c.ShouldBind(&form); err != nil {
			applyError(&entry, tools, c, "Binding Form", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": entry.Response})
			return
		}

		if !middleware.VerifyStruct(&form) {
			applyError(&entry, tools, c, "Struct Verify", errors.New("cannot verify struct"))
			c.JSON(http.StatusBadRequest, gin.H{"error": entry.Response})
			return
		}

		if !middleware.VerifyEmail(&form, tools) {
			applyError(&entry, tools, c, "Struct Verify", errors.New("cannot verify struct"))
			c.JSON(http.StatusBadRequest, gin.H{"error": entry.Response})
			return
		}

		entry.Name = form.Name
		entry.Email = form.Email
		entry.Subject = form.Subject
		entry.Comment = form.Comment

		success, err := middleware.VerifyTurnstile(form.CFTurnstileResponse, tools.Client)
		if err != nil {
			applyError(&entry, tools, c, "Turnstile Verification", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": entry.Response})
			return
		} else if !success {
			applyError(&entry, tools, c, "Turnstile Verification", errors.New("failed turnstile challenge"))
			c.JSON(http.StatusBadRequest, gin.H{"error": entry.Response})
			return
		}

		ct, err := database.GetCountOfSuccessfulEntries(tools.DB)
		if err != nil {
			applyError(&entry, tools, c, "Query Entry Count", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": entry.Response})
			return
		} else if ct >= 50 {
			applyError(&entry, tools, c, "Query Entry After Count", errors.New("at/past capacity: "+strconv.Itoa(ct)))
			c.JSON(http.StatusBadRequest, gin.H{"error": entry.Response})
			return
		}

		if err := SendConfirmationSlack(tools.Client, entry); err != nil {
			applyError(&entry, tools, c, "Slack error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": entry.Response})
			return
		}

		entry.Complete = true

		if err := SendConfirmationToUser(tools.SendGrid, entry.Name, entry.Email); err != nil {
			applyError(&entry, tools, c, "Message To User", err)
			c.Status(204)
			return
		}

		if err := SendConfirmationEmail(tools.SendGrid, entry); err != nil {
			applyError(&entry, tools, c, "Message To Myself", err)
			c.Status(204)
			return
		}

		if err := database.InsertEntry(tools.DB, entry); err != nil {
			fmt.Printf("Unable to save entry: %+v\n", entry)
		}
		c.Status(204)

	}
}
