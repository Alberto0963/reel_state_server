package models

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// // Estructura para los datos de la tarjeta y pago
// type CardDetails struct {
//     HolderName string `json:"holder_name"`
//     CardNumber string `json:"card_number"`
//     Cvv2       string `json:"cvv2"`
//     ExpMonth   string `json:"expiration_month"`
//     ExpYear    string `json:"expiration_year"`
// }

type SubscriptionResponse struct {
	ID                string `json:"id"`
	CreationDate      string `json:"creation_date"`
	Status            string `json:"status"`
	ChargeDate        string `json:"charge_date"`
	CancelAtPeriodEnd bool   `json:"cancel_at_period_end"`
	TrialEndDate      string `json:"trial_end_date"`
	PeriodEndDate     string `json:"period_end_date"`
	PlanID            string `json:"plan_id"`
	CustomerID        string `json:"customer_id"`
	// Card                Card    `json:"card"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type ChargeRequest struct {
	SourceID string `json:"source_id"`
	Method   string `json:"method"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	MembershipID    string `json:"membership_id"`
	PlanID string `json:"plan_id"`

	Amount        float64 `json:"amount"`
	Description   string  `json:"description"`
	DeviceSession string  `json:"device_session_id"`
	CustomerID    string  `json:"CUSTOMER_ID"`
}

type Customer struct {
	Name        string `json:"name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	// Address    Address `json:"address"`
}

// Función para realizar el cargo
func MakePayment(token string, amount float64, description string, deviceSessionId string, customer string) error {
	url := os.Getenv("API_URL") + os.Getenv("MERCHANT_ID") + "/customers/" + customer + "/subscriptions"

	chargeReq := ChargeRequest{
		// SourceID:      token,
		Method:        "card",
		Amount:        amount,
		Description:   description,
		DeviceSession: deviceSessionId,
	}

	chargeData, err := json.Marshal(chargeReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(chargeData))
	if err != nil {
		return err
	}

	// Autenticación básica con la llave privada de Openpay
	req.Header.Set("Authorization", "Basic "+os.Getenv("OPENPAY_PRIVATE_KEY"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return fmt.Errorf("Error en el pago: %s", string(body))
	}

	fmt.Println("Pago realizado con éxito:", string(body))
	return nil
}

type SubscriptionRequest struct {
	PlanID    string  `json:"plan_id"`
	SourceID  string  `json:"source_id"`
	TrialDays int     `json:"trial_days,omitempty"`
	Amount    float64 `json:"amount,omitempty"`
	Currency  string  `json:"currency,omitempty"`
}

func CreateSubscriptionWithToken(customerID, planID, tokenID string, amount float64, trialDays int) (SubscriptionResponse, error) {
	merchantID := os.Getenv("MERCHANT_ID")
	apiKey := os.Getenv("OPENPAY_PRIVATE_KEY")
	var sub *SubscriptionResponse

	subscriptionRequest := SubscriptionRequest{
		PlanID:    planID,
		SourceID:  tokenID,
		TrialDays: trialDays,
		Amount:    amount,
		Currency:  "MXN",
	}

	body, err := json.Marshal(subscriptionRequest)
	if err != nil {
		return *sub, fmt.Errorf("error al generar la solicitud de suscripción: %v", err)
	}

	url := fmt.Sprintf("https://sandbox-api.openpay.mx/v1/%s/customers/%s/subscriptions", merchantID, customerID)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return *sub, fmt.Errorf("error al crear la solicitud HTTP: %v", err)
	}

	auth := base64.StdEncoding.EncodeToString([]byte(apiKey + ":"))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return *sub, fmt.Errorf("error al realizar la solicitud HTTP: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		return *sub, fmt.Errorf("error al crear suscripción: %v", errorResponse)
	}
	// Parseamos la respuesta de la suscripción
	sub, err = ParseSubscriptionResponse(resp)
	if err != nil {
		return *sub,	fmt.Errorf("Error al parsear la respuesta de la suscripción: %v\n", err)

	}

	fmt.Println("Suscripción creada con éxito")
	return *sub, nil
}

func CreateCustomer(user Customer) (string, error) {
	merchantID := os.Getenv("MERCHANT_ID")     // Merchant ID de Openpay
	apiKey := os.Getenv("OPENPAY_PRIVATE_KEY") // Clave privada de Openpay
	urlapi := os.Getenv("API_URL")
	customer := Customer{
		Name:        user.Name,
		LastName:    user.LastName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		// Address: Address{
		// 	Line1:      "Calle Morelos #12",
		// 	PostalCode: "76000",
		// 	State:      "Querétaro",
		// 	City:       "Querétaro",
		// 	CountryCode: "MX",
		// },
	}

	body, err := json.Marshal(customer)
	if err != nil {
		return "", fmt.Errorf("error al generar la solicitud de cliente: %v", err)
	}

	url := fmt.Sprintf("%s/%s/customers", urlapi, merchantID)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("error al crear la solicitud HTTP: %v", err)
	}

	auth := base64.StdEncoding.EncodeToString([]byte(apiKey + ":"))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error al realizar la solicitud HTTP: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		return "", fmt.Errorf("error al crear cliente: %v", errorResponse)
	}

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	customerID := response["id"].(string)

	fmt.Println("Cliente creado con éxito. ID:", customerID)
	return customerID, nil
}

// Función que recibe el cuerpo de la respuesta HTTP y lo mapea al modelo Subscription
func ParseSubscriptionResponse(resp *http.Response) (*SubscriptionResponse, error) {
	// Leemos el cuerpo de la respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer el cuerpo de la respuesta: %v", err)
	}
	defer resp.Body.Close()

	// Creamos una instancia del modelo Subscription
	var subscription SubscriptionResponse

	// Convertimos el JSON al modelo Subscription
	err = json.Unmarshal(body, &subscription)
	if err != nil {
		return nil, fmt.Errorf("error al decodificar el JSON de la suscripción: %v", err)
	}

	return &subscription, nil
}

func CancelOpenPaySubscription(customerID string,subscriptionID string)error {
	// Obtener los parámetros necesarios
	// customerID := c.Param("customerID")
	// subscriptionID := c.Param("subscriptionID")
	// merchantID := os.Getenv("MERCHANT_ID")     // Merchant ID de Openpay
	// apiKey := os.Getenv("OPENPAY_PRIVATE_KEY") // Clave privada de Openpay
	
	urlapi := os.Getenv("API_URL")

	// Obtener las credenciales de Openpay desde las variables de entorno
	merchantID := os.Getenv("MERCHANT_ID")
	apiKey := os.Getenv("OPENPAY_PRIVATE_KEY")

	// Construir la URL para cancelar la suscripción
	url := fmt.Sprintf("%s%s/customers/%s/subscriptions/%s",urlapi, merchantID, customerID, subscriptionID)

	// Crear la solicitud HTTP DELETE
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear la solicitud HTTP"})
		return err
	}

	// Agregar la autenticación básica
	auth := base64.StdEncoding.EncodeToString([]byte(apiKey + ":"))
	req.Header.Set("Authorization", "Basic "+auth)

	// Enviar la solicitud
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al realizar la solicitud HTTP"})
		return err
	}
	defer resp.Body.Close()

	// Verificar la respuesta
	if resp.StatusCode == http.StatusNoContent {
		return nil
		// c.JSON(http.StatusOK, gin.H{"message": "Suscripción cancelada exitosamente"})
	} else {
		var errorResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			// c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar la respuesta del servidor"})
			return err
		}
		return err
		// c.JSON(resp.StatusCode, gin.H{"error": errorResponse})
	}
}
