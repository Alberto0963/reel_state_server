package controllers

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"github.com/gin-gonic/gin"
// )

// // Estructura para la verificación del webhook
// type WebhookVerification struct {
// 	Mode      string `json:"hub.mode"`
// 	Token     string `json:"hub.verify_token"`
// 	Challenge string `json:"hub.challenge"`
// }

// // Estructura para manejar eventos de WhatsApp
// type WhatsAppMessage struct {
// 	Entry []struct {
// 		Changes []struct {
// 			Value struct {
// 				Messages []struct {
// 					From    string `json:"from"`
// 					ID      string `json:"id"`
// 					Text    struct {
// 						Body string `json:"body"`
// 					} `json:"text"`
// 				} `json:"messages"`
// 			} `json:"value"`
// 		} `json:"changes"`
// 	} `json:"entry"`
// }

// // Webhook para WhatsApp con Gin
// func WebhookHandlerWhatsapp(c *gin.Context) {
// 	verifyToken := os.Getenv("Whatsapp_Token")

// 	switch c.Request.Method {
// 	case http.MethodGet:
// 		// Verificación del webhook con Meta
// 		mode := c.Query("hub.mode")
// 		token := c.Query("hub.verify_token")
// 		challenge := c.Query("hub.challenge")

// 		if mode == "subscribe" && token == verifyToken {
// 			c.String(http.StatusOK, challenge)
// 			return
// 		}

// 		c.JSON(http.StatusForbidden, gin.H{"error": "No autorizado"})

// 	case http.MethodPost:
// 		var message WhatsAppMessage

// 		// Leer y parsear el JSON del body
// 		if err := c.ShouldBindJSON(&message); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Error en formato JSON"})
// 			return
// 		}

// 		// Procesar el mensaje
// 		for _, entry := range message.Entry {
// 			for _, change := range entry.Changes {
// 				for _, msg := range change.Value.Messages {
// 					fmt.Printf("Mensaje recibido de %s: %s\n", msg.From, msg.Text.Body)
// 				}
// 			}
// 		}

// 		c.Status(http.StatusOK)
// 	}
// }
import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reelState/services"

	// "io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// Estructura para recibir mensajes de WhatsApp
type WhatsAppMessage struct {
	Entry []struct {
		Changes []struct {
			Value struct {
				Messages []struct {
					From    string `json:"from"`
					Text    struct {
						Body string `json:"body"`
					} `json:"text"`
				} `json:"messages"`
			} `json:"value"`
		} `json:"changes"`
	} `json:"entry"`
}

// Webhook para WhatsApp
func WebhookHandlerWhatsapp(c *gin.Context) {
	verifyToken := os.Getenv("WHATSAPP_TOKEN")

	if c.Request.Method == http.MethodGet {
		// Verificar el webhook con Meta
		if c.Query("hub.verify_token") == verifyToken {
			c.String(http.StatusOK, c.Query("hub.challenge"))
			return
		}
		c.JSON(http.StatusForbidden, gin.H{"error": "Token inválido"})
		return
	}

	// Procesar mensajes entrantes
	var message WhatsAppMessage
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	// Obtener mensaje y número
	for _, entry := range message.Entry {
		for _, change := range entry.Changes {
			for _, msg := range change.Value.Messages {
				textoRecibido := msg.Text.Body
				numero := msg.From

				// Enviar mensaje a Vertex AI
				respuestaIA,err := services.GenerateResponse(textoRecibido)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				// GenerarRespuestaVertexAI(textoRecibido)

				// Responder en WhatsApp
				EnviarMensajeWhatsApp(numero, respuestaIA)
			}
		}
	}

	c.Status(http.StatusOK)
}


// Función para enviar mensajes a WhatsApp
// EnviarMensajeWhatsApp envía un mensaje por WhatsApp y maneja errores
func EnviarMensajeWhatsApp(numero string, mensaje string) error {
	token := os.Getenv("Whatsapp_Token")
	phoneID := os.Getenv("WHATSAPP_PHONE_ID")

	if token == "" || phoneID == "" {
		return errors.New("falta configurar WHATSAPP_TOKEN o WHATSAPP_PHONE_ID en el entorno")
	}

	url := "https://graph.facebook.com/v17.0/" + phoneID + "/messages"

	// Crear JSON del mensaje
	reqBody, err := json.Marshal(map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                numero,
		"type":              "text",
		"text": map[string]string{
			"body": mensaje,
		},
	})
	if err != nil {
		return fmt.Errorf("error al generar JSON del mensaje: %v", err)
	}

	// Crear solicitud HTTP
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("error al crear la solicitud HTTP: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Enviar la solicitud
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error en la petición HTTP: %v", err)
	}
	defer resp.Body.Close()

	// Verificar el código de respuesta
	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		return fmt.Errorf("error en la respuesta de WhatsApp: %v", errorResp)
	}

	fmt.Println("Mensaje enviado a", numero)
	return nil
}

