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

	case "PAYMENT.SALE.COMPLETED":
		// Obtener la fecha del pago
		paymentDateStr :=  event.CreateTime
		paymentDate, _ := time.Parse(time.RFC3339, paymentDateStr)

		// Calcular la siguiente fecha de pago (suma 30 días)
		nextPaymentDate := paymentDate.AddDate(0, 1, 0)

		// Actualizar la base de datos
		models.UpdateSubscriptionWebhook(event.Resource.ID,paymentDate, nextPaymentDate, true)

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
		
		models.CancelSubscriptionWebhook(event.Resource.ID, parsedTime)

	default:
		fmt.Printf("Unhandled event type: %s\n", event.EventType)
	}

	// Respond to the webhook event
	c.JSON(http.StatusOK, gin.H{"message": "Event received"})
}

// OpenpayWebhook representa el webhook enviado por Openpay.
type OpenpayWebhook struct {
	Type             string      `json:"type"`
	EventDate        time.Time   `json:"event_date"`
	Transaction      Transaction `json:"transaction"`
	VerificationCode string      `json:"verification_code"`
}

// Transaction representa la transacción en el webhook de Openpay.
type Transaction struct {
	ID              string    `json:"id"`
	Authorization   string    `json:"authorization"`
	OperationType   string    `json:"operation_type"`
	TransactionType string    `json:"transaction_type"`
	Status          string    `json:"status"`
	Conciliated     bool      `json:"conciliated"`
	CreationDate    time.Time `json:"creation_date"`
	OperationDate   time.Time `json:"operation_date"`
	Description     string    `json:"description"`
	ErrorMessage    string    `json:"error_message,omitempty"`
	OrderID         string    `json:"order_id,omitempty"`
	Card            Card      `json:"card"`
	CustomerID      string    `json:"customer_id"`
	SubscriptionID  string    `json:"subscription_id"`
	Fee             Fee       `json:"fee"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	Method          string    `json:"method"`
}

// Card representa la información de la tarjeta en la transacción.
type Card struct {
	Type            string `json:"type"`
	Brand           string `json:"brand"`
	Address         string `json:"address,omitempty"`
	CardNumber      string `json:"card_number"`
	HolderName      string `json:"holder_name"`
	ExpirationYear  string `json:"expiration_year"`
	ExpirationMonth string `json:"expiration_month"`
	AllowsCharges   bool   `json:"allows_charges"`
	AllowsPayouts   bool   `json:"allows_payouts"`
	BankName        string `json:"bank_name"`
	BankCode        string `json:"bank_code"`
}

// Fee representa la comisión aplicada en la transacción.
type Fee struct {
	Amount         float64  `json:"amount"`
	Tax            float64  `json:"tax"`
	Surcharge      *float64 `json:"surcharge,omitempty"`
	BaseCommission *float64 `json:"base_commission,omitempty"`
	Currency       string   `json:"currency"`
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
			fmt.Printf("Codigo de Verificacion: %s", webhook.VerificationCode)

			c.String(http.StatusOK, verificationCode)
			return
		}
	}
	// fmt.Printf("Codigo de Verificacion: %s", webhook.EventType)

	// Manejar diferentes tipos de eventos
	switch webhook.Type {
	case "verification":
		err := handleVerificationCode(webhook)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"message": "éxito"})

		}

	case "subscription.create", "charge.succeeded":
		err := handleSubscriptionCreated(webhook)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"message": "éxito"})

		}
	case "charge.refund":
		err := handlePaymentRefunded(webhook)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"message": "Reembolso realizado"})

		}
	case "charge.failed":
		err := handlePaymentDeclined(webhook)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"message": "Pago declinado"})

		}
	case "subscription.cancel":
		err := handleSubscriptionCanceled(webhook)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"message": "Reembolso realizado"})

		}
	default:
		c.JSON(http.StatusOK, gin.H{"message": "Evento no manejado"})
	}
}

// Función para manejar la creación de suscripciones
func handleVerificationCode(webhook OpenpayWebhook) error {
	vc := webhook.VerificationCode

	// Realizar las operaciones necesarias (guardar en la base de datos, enviar confirmación, etc.)
	fmt.Printf("Codigo de Verificacion: %s", vc)
	return nil
	// c.JSON(http.StatusOK, gin.H{"message": "Exito", "codigo verificaion": vc})
}

// Función para manejar la creación de suscripciones
func handleSubscriptionCreated(webhook OpenpayWebhook) error {
	subscriptionID := webhook.Transaction.SubscriptionID
	// subscriptionID := webhook.Transaction.SubscriptionID
	status := webhook.Transaction.Status

	// Usar directamente EventDate como time.Time
	paymentDate := webhook.EventDate

	// Calcular la siguiente fecha de pago (suma 1 mes)
	nextPaymentDate := paymentDate.AddDate(0, 1, 0)

	// Actualizar la base de datos con las fechas de pago
	models.UpdateSubscriptionWebhook(webhook.Transaction.SubscriptionID, paymentDate, nextPaymentDate, true)
	// Realizar las operaciones necesarias (guardar en la base de datos, enviar confirmación, etc.)
	fmt.Printf("Suscripción creada:\nID: %s \nEstado: %s\n", subscriptionID, status)

	return nil
	// c.JSON(http.StatusOK, gin.H{"message": "Suscripción creada con éxito", "subscription_id": subscriptionID})
}

// Función para manejar pagos declinados
func handlePaymentDeclined(webhook OpenpayWebhook) error {
	transactionID := webhook.Transaction.ID
	errorMessage := webhook.Transaction.ErrorMessage

	// Realizar las operaciones necesarias (enviar notificación, actualizar el estado, etc.)
	fmt.Printf("Pago declinado:\nID: %s\nMensaje de error: %s\n", transactionID, errorMessage)
	// layout := "2006-01-02T15:04:05Z"

	// Parse the date string
	// parsedTime, err := time.Parse(layout, webhook.EventDate)
	// if err != nil {
	// 	fmt.Println("Error parsing date:", err)
	// 	return err
	// }

	err := models.CancelSubscriptionWebhook(webhook.Transaction.SubscriptionID, webhook.EventDate)
	if err != nil {
		fmt.Println("Error:", err)
		// c.JSON(http.StatusOK, gin.H{"message": "Pago declinado", "transaction_id": transactionID, "error": errorMessage})

		return err
	}
	return nil
	// c.JSON(http.StatusOK, gin.H{"message": "Pago declinado", "transaction_id": transactionID, "error": errorMessage})
}

// Función para manejar la cancelación de suscripciones
func handleSubscriptionCanceled(webhook OpenpayWebhook) error {
	subscriptionID := webhook.Transaction.SubscriptionID
	// planID := webhook.Subscription.PlanID
	status := webhook.Transaction.Status

	// Realizar las operaciones necesarias (actualizar la suscripción en la base de datos, etc.)
	fmt.Printf("Suscripción cancelada:\nID: %s\nEstado: %s\n", subscriptionID, status)
	return nil
	// c.JSON(http.StatusOK, gin.H{"message": "Suscripción cancelada", "subscription_id": subscriptionID})
}

// Función para manejar reembolsos de pagos (opcional)
func handlePaymentRefunded(webhook OpenpayWebhook) error {
	transactionID := webhook.Transaction.ID
	amount := webhook.Transaction.Amount

	// Realizar las operaciones necesarias (actualizar la base de datos, etc.)
	fmt.Printf("Reembolso realizado:\nID: %s\nMonto: %f\n", transactionID, amount)
	return nil
	// c.JSON(http.StatusOK, gin.H{"message": "Reembolso realizado", "transaction_id": transactionID, "amount": amount})
}
