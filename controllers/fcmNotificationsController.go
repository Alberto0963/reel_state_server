package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reelState/models"
	"reelState/utils/token"

	// firebase "firebase.google.com/go"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/gin-gonic/gin"

	// "google.golang.org/api/fcm/v1"
	"google.golang.org/api/option"
)

// Estructura para la solicitud FCM
type FCMRequest struct {
	To           string                 `json:"to"`
	Data         map[string]interface{} `json:"data"`
	Notification map[string]string      `json:"notification"`
}

// Estructura para la solicitud Huawei
type HuaweiRequest struct {
	ValidateOnly bool          `json:"validate_only"`
	Message      HuaweiMessage `json:"message"`
}

type HuaweiMessage struct {
	Token        []string           `json:"token"`
	Notification HuaweiNotification `json:"notification"`
}

type HuaweiNotification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// Enviar notificación a través de FCM o Huawei
func SendNotification(token, title, body, deviceType string) error {
	switch deviceType {
	case "android", "ios":
		return sendFCMNotification(token, title, body)
	case "huawei":
		return sendHuaweiNotification(token, title, body)
	default:
		return fmt.Errorf("Tipo de dispositivo no compatible")
	}
}

// // Enviar notificación a FCM
// func sendFCMNotification(token, title, body string) error {
//     fcmURL := "https://fcm.googleapis.com/fcm/send"
//     // serverKey := "YOUR_FCM_SERVER_KEY"
// 	serverKey := os.Getenv("FCM_SERVER_KEY")

//     requestBody := FCMRequest{
//         To: token,
//         Notification: map[string]string{
//             "title": title,
//             "body":  body,
//         },
//         Data: map[string]interface{}{
//             "click_action": "FLUTTER_NOTIFICATION_CLICK",
//         },
//     }

//     jsonBody, err := json.Marshal(requestBody)
//     if err != nil {
//         return err
//     }

//     req, err := http.NewRequest("POST", fcmURL, bytes.NewBuffer(jsonBody))
//     if err != nil {
//         return err
//     }

//     req.Header.Set("Authorization", "key="+serverKey)
//     req.Header.Set("Content-Type", "application/json")

//     client := &http.Client{}
//     res, err := client.Do(req)
//     if err != nil {
//         return err
//     }
//     defer res.Body.Close()

//     if res.StatusCode != http.StatusOK {
//         return fmt.Errorf("Error al enviar notificación FCM")
//     }

//     return nil
// }

// func sendFCMNotification(token string, title string, body string) error {
// 	ctx := context.Background()

// 	// Cargar las credenciales desde el archivo JSON
// 	serviceAccountFile := "/home/alberto/Downloads/reelstate-8cc46-firebase.json"                                          // Cambia esto a la ruta de tu archivo JSON
// 	fcmService, err := fcm.NewService(ctx, option.WithCredentialsFile(serviceAccountFile)) // Cambia esto a la ruta de tu archivo JSON
// 	if err != nil {
// 		return fmt.Errorf("no se pudo crear el token: %v", err)
// 	}

// 	// Crear el cliente FCM
// 	// fcmService, err = fcm.NewService(ctx)
// 	// if err != nil {
// 	// 	return fmt.Errorf("no se pudo crear el servicio FCM: %v", err)
// 	// }

// 	// Preparar la notificación
// 	message := &fcm.Message{
// 		Token: token,
// 		Notification: &fcm.Notification{
// 			Title: title,
// 			Body:  body,
// 		},
// 	}

// 	// Crear la solicitud de envío
// 	request := &fcm.SendMessageRequest{
// 		Message: message,
// 	}
// 	// Enviar la notificación
// 	response, err := fcmService.Projects.Messages.Send("projects/reelstate-8cc46/messages:send", request).Do()
// 	if err != nil {
// 		return fmt.Errorf("error al enviar la notificación: %v", err)
// 	}

// 	log.Printf("Respuesta de FCM: %v", response)
// 	return nil
// }

func InitFirebaseApp() (*firebase.App, error) {
	// Cargar las credenciales del archivo JSON
	serverKey := os.Getenv("FCM_SERVER_KEY")

	opt := option.WithCredentialsFile(serverKey)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func sendFCMNotification(token, title, body string) error {
	// Inicializar la app de Firebase
	app, err := InitFirebaseApp()
	if err != nil {
		return err
	}

	// Crear un cliente de mensajería
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		return err
	}

	// Crear el mensaje
	message := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: map[string]string{
			"key": "value",
		},
	}

	// Enviar la notificación
	response, err := client.Send(ctx, message)
	if err != nil {
		return err
	}

	log.Printf("Mensaje enviado exitosamente: %s\n", response)
	return nil
}

// Enviar notificación a Huawei Push Kit
func sendHuaweiNotification(token, title, body string) error {
	huaweiURL := "https://push-api.cloud.huawei.com/v1/YOUR_APP_ID/messages:send"
	accessToken := "YOUR_HUAWEI_ACCESS_TOKEN"

	requestBody := HuaweiRequest{
		ValidateOnly: false,
		Message: HuaweiMessage{
			Token: []string{token},
			Notification: HuaweiNotification{
				Title: title,
				Body:  body,
			},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", huaweiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Error al enviar notificación a Huawei")
	}

	return nil
}

// Estructura de la solicitud para actualizar el token
type UpdateTokenRequest struct {
	TokenDevice string `json:"token_device" binding:"required"`
}

// Función para actualizar el token del dispositivo
func UpdateDeviceToken(c *gin.Context) {
	// Obtener el ID del usuario desde el token
	userID, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	// Parsear el cuerpo de la solicitud
	var req UpdateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token del dispositivo requerido"})
		return
	}

	// Buscar el usuario y actualizar el token del dispositivo
	user, err := models.GetUserByIdToUpdate(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el token"})
		return
	}
	user.DeviceToken = req.TokenDevice
	user.UpdateUser()
	// models.UserUpdate(user)
	// Responder con éxito
	c.JSON(http.StatusOK, gin.H{"message": "Token del dispositivo actualizado correctamente"})
}
