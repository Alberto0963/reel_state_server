package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reelState/models"
	"time"

	"github.com/gin-gonic/gin"
)

// HandleWebhook handles the webhook request
func HandleWebhook(c *gin.Context) {
	// Read the body of the request
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to read request body"})
		return
	}

	// Decode the webhook event
	var event models.WebhookEvent
	err = json.Unmarshal(body, &event)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse request body"})
		return
	}

	// Handle the event based on its type
	switch event.EventType {
	case "BILLING.SUBSCRIPTION.CREATED":
		// Process the billing subscription created event
		fmt.Printf("Billing subscription created: %+v\n", event)
		
	case "BILLING.SUBSCRIPTION.PAYMENT.FAILED", "BILLING.SUBSCRIPTION.CANCELLED":

		// var sub models.Subscription

		// sub.PaypalSubscriptionID =
		// sub.RenewalCancelledAt =
		// Define the layout (format) of the date string
		layout := "2006-01-02T15:04:05Z"

		// Parse the date string
		parsedTime, err := time.Parse(layout, event.CreateTime)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			return
		}
		models.CancelSubscriptionFunction(event.Resource.ID, parsedTime)

	default:
		fmt.Printf("Unhandled event type: %s\n", event.EventType)
	}

	// Respond to the webhook event
	c.JSON(http.StatusOK, gin.H{"message": "Event received"})
}
