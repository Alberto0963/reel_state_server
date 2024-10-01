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

// Estructura para recibir la información del webhook
type OpenpayWebhook struct {
	EventType        string `json:"event_type"`
	EventDate        string `json:"event_date"`
	VerificationCode string `json:"verification_code"`
	Transaction      struct {
		ID           string  `json:"id"`
		Status       string  `json:"status"`
		Amount       float64 `json:"amount,omitempty"`
		PlanID       string  `json:"plan_id,omitempty"`
		ErrorMessage string  `json:"error_message,omitempty"`
	} `json:"transaction"`
	Subscription struct {
		ID     string `json:"id"`
		Status string `json:"status"`
		PlanID string `json:"plan_id,omitempty"`
	} `json:"subscription,omitempty"`
}

// Webhook para manejar diferentes tipos de eventos de Openpay
func OpenpayWebhookHandler(c *gin.Context) {
	var webhook OpenpayWebhook

	// Decodificar el cuerpo del webhook
	if err := c.ShouldBindJSON(&webhook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error al procesar el webhook", "error": err.Error()})
		return
	}

	// Verificar si es una solicitud de verificación de webhook
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err == nil {
		if verificationCode, ok := payload["verification_code"].(string); ok {
			// Responder con el código de verificación
			c.String(http.StatusOK, verificationCode)
			return
		}
	}

	// Manejar diferentes tipos de eventos
	switch webhook.EventType {
	case "verification":
		handleVerificationCode(webhook, c)
	case "subscription.create":
		handleSubscriptionCreated(webhook, c)
	case "charge.refund":
		handlePaymentRefunded(webhook, c)
	case "charge.failed":
		handlePaymentDeclined(webhook, c)
	case "subscription.cancel":
		handleSubscriptionCanceled(webhook, c)
	default:
		c.JSON(http.StatusOK, gin.H{"message": "Evento no manejado"})
	}
}

// Función para manejar la creación de suscripciones
func handleVerificationCode(webhook OpenpayWebhook, c *gin.Context) {
	vc := webhook.VerificationCode


	// Realizar las operaciones necesarias (guardar en la base de datos, enviar confirmación, etc.)
	fmt.Printf("Codigo de Verificacion: %s", vc)

	c.JSON(http.StatusOK, gin.H{"message": "Exito", "codigo verificaion": vc})
}

// Función para manejar la creación de suscripciones
func handleSubscriptionCreated(webhook OpenpayWebhook, c *gin.Context) {
	subscriptionID := webhook.Subscription.ID
	planID := webhook.Subscription.PlanID
	status := webhook.Subscription.Status

	// Realizar las operaciones necesarias (guardar en la base de datos, enviar confirmación, etc.)
	fmt.Printf("Suscripción creada:\nID: %s\nPlanID: %s\nEstado: %s\n", subscriptionID, planID, status)

	c.JSON(http.StatusOK, gin.H{"message": "Suscripción creada con éxito", "subscription_id": subscriptionID})
}

// Función para manejar pagos declinados
func handlePaymentDeclined(webhook OpenpayWebhook, c *gin.Context) {
	transactionID := webhook.Transaction.ID
	errorMessage := webhook.Transaction.ErrorMessage

	// Realizar las operaciones necesarias (enviar notificación, actualizar el estado, etc.)
	fmt.Printf("Pago declinado:\nID: %s\nMensaje de error: %s\n", transactionID, errorMessage)
	layout := "2006-01-02T15:04:05Z"

	// Parse the date string
	parsedTime, err := time.Parse(layout, webhook.EventDate)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}

	err = models.CancelSubscriptionFunction(webhook.Subscription.ID, parsedTime)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		c.JSON(http.StatusOK, gin.H{"message": "Pago declinado", "transaction_id": transactionID, "error": errorMessage})

		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Pago declinado", "transaction_id": transactionID, "error": errorMessage})
}

// Función para manejar la cancelación de suscripciones
func handleSubscriptionCanceled(webhook OpenpayWebhook, c *gin.Context) {
	subscriptionID := webhook.Subscription.ID
	planID := webhook.Subscription.PlanID
	status := webhook.Subscription.Status

	// Realizar las operaciones necesarias (actualizar la suscripción en la base de datos, etc.)
	fmt.Printf("Suscripción cancelada:\nID: %s\nPlanID: %s\nEstado: %s\n", subscriptionID, planID, status)

	c.JSON(http.StatusOK, gin.H{"message": "Suscripción cancelada", "subscription_id": subscriptionID})
}

// Función para manejar reembolsos de pagos (opcional)
func handlePaymentRefunded(webhook OpenpayWebhook, c *gin.Context) {
	transactionID := webhook.Transaction.ID
	amount := webhook.Transaction.Amount

	// Realizar las operaciones necesarias (actualizar la base de datos, etc.)
	fmt.Printf("Reembolso realizado:\nID: %s\nMonto: %f\n", transactionID, amount)

	c.JSON(http.StatusOK, gin.H{"message": "Reembolso realizado", "transaction_id": transactionID, "amount": amount})
}
