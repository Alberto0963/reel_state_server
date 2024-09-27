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

type ChargeRequest struct {
	SourceID string `json:"source_id"`
	Method   string `json:"method"`
	Name     string `json:"name"`
    Email     string `json:"email"`

	PlanID   string `json:"plan_id"`

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

func CreateSubscriptionWithToken(customerID, planID, tokenID string, amount float64, trialDays int) error {
	merchantID := os.Getenv("MERCHANT_ID")
	apiKey := os.Getenv("OPENPAY_PRIVATE_KEY")

	subscriptionRequest := SubscriptionRequest{
		PlanID:    planID,
		SourceID:  tokenID,
		TrialDays: trialDays,
		Amount:    amount,
		Currency:  "MXN",
	}

	body, err := json.Marshal(subscriptionRequest)
	if err != nil {
		return fmt.Errorf("error al generar la solicitud de suscripción: %v", err)
	}

	url := fmt.Sprintf("https://sandbox-api.openpay.mx/v1/%s/customers/%s/subscriptions", merchantID, customerID)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error al crear la solicitud HTTP: %v", err)
	}

	auth := base64.StdEncoding.EncodeToString([]byte(apiKey + ":"))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error al realizar la solicitud HTTP: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		return fmt.Errorf("error al crear suscripción: %v", errorResponse)
	}

	fmt.Println("Suscripción creada con éxito")
	return nil
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

	url := fmt.Sprintf("%s/%s/customers",urlapi, merchantID)

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
