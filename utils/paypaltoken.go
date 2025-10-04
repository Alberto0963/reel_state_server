package SMS

import (
	"encoding/json"
	"fmt"
	// "log"
	"net/http"
	"net/url"
	// "os"
	"strings"

	// "time"

	// "github.com/joho/godotenv"
)

// Estructura para el token de respuesta
type TokenResponse struct {
    AccessToken string `json:"access_token"`
    ExpiresIn   int    `json:"expires_in"`
}

// Función para obtener un nuevo token desde PayPal
func FetchPayPalToken(clientID, secret string) (string, error) {
    // Endpoint de PayPal para obtener el token
    tokenURL := "https://api.paypal.com/v1/oauth2/token"

    // Cuerpo de la solicitud
    data := url.Values{}
    data.Set("grant_type", "client_credentials")

    // Crear una nueva solicitud HTTP
    req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
    if err != nil {
        return "", err
    }

    // Autenticación con Basic Auth
    req.SetBasicAuth(clientID, secret)
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    // Enviar la solicitud HTTP
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    // Verificar el código de respuesta HTTP
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("failed to fetch token: status code %d", resp.StatusCode)
    }

    // Decodificar la respuesta JSON para obtener el token
    var tokenResponse TokenResponse
    if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
        return "", err
    }

    // Retornar el token de acceso
    return tokenResponse.AccessToken, nil
}




