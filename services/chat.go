package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	// "io/ioutil"
	// "log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	// "google.golang.org/api/idtoken"
	// "google.golang.org/api/idtoken"
)

// Estructura para la solicitud a Vertex AI
type VertexAIRequest struct {
	Model             string            `json:"model"`
	Contents          []Content         `json:"contents"`
	Config            Config            `json:"generation_config"`
	SystemInstruction SystemInstruction `json:"system_instruction"` // Cambio: De []Part a SystemInstruction
	SafetySettings    []SafetySetting   `json:"safety_settings"`
}

type Content struct {
	Role  string `json:"role"`
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type Config struct {
	Temperature        float64  `json:"temperature"`
	TopP               float64  `json:"top_p"`
	MaxOutputTokens    int      `json:"max_output_tokens"`
	ResponseModalities []string `json:"response_modalities"` // Falta en tu struct
}

// Estructura para system_instruction (antes era []Part)
type SystemInstruction struct {
	Role  string `json:"role"`
	Parts []Part `json:"parts"`
}

type SafetySetting struct {
	Category  string `json:"category"`
	Threshold string `json:"threshold"`
}

// ... (Structs VertexAIRequest, Content, Part, Config, SafetySetting sin cambios)

// func GenerateResponse(userMessage string) (string, error) {
// 	ctx := context.Background()

// 	// 1. Obtener las credenciales de Google Cloud (reemplaza con tu método preferido)
// 	// a. Credenciales desde archivo JSON (recomendado para desarrollo local):
// 	// credentialsFile := "path/to/your/credentials.json" // Ruta al archivo JSON
// 	// creds, err := google.CredentialsFromFile(ctx, credentialsFile, "https://www.googleapis.com/auth/cloud-platform")
// 	// if err != nil {
// 	//     return "", fmt.Errorf("error obteniendo credenciales: %w", err)
// 	// }

// 	// b. Credenciales desde variables de entorno (para despliegue en Google Cloud):
// 	creds, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/cloud-platform")
// 	if err != nil {
// 			return "", fmt.Errorf("error obteniendo credenciales: %w", err)
// 	}

// 	// 2. Crear el cliente HTTP con las credenciales
// 	client := oauth2.NewClient(ctx, creds.TokenSource) // Usa creds.TokenSource

// 	url := "https://us-central1-aiplatform.googleapis.com/v1/projects/reelstate-8cc46/locations/us-central1/publishers/google/models/gemini-2.0:generateContent"

// 	requestBody := VertexAIRequest{
// 		Model: "gemini-2.0",
// 		Contents: []Content{
// 			{Role: "user", Parts: []Part{{Text: userMessage}}},
// 		},
// 		Config: Config{
// 			Temperature:      1.0,
// 			TopP:            0.95,
// 			MaxOutputTokens:  8192,
// 			SafetySettings: []SafetySetting{
// 				{Category: "HARM_CATEGORY_HATE_SPEECH", Threshold: "OFF"},
// 				{Category: "HARM_CATEGORY_DANGEROUS_CONTENT", Threshold: "OFF"},
// 				{Category: "HARM_CATEGORY_SEXUALLY_EXPLICIT", Threshold: "OFF"},
// 				{Category: "HARM_CATEGORY_HARASSMENT", Threshold: "OFF"},
// 			},
// 			SystemInstruction: []Part{
// 				{Text: `Eres un chatbot especializado en bienes raíces...`}, // Instrucciones completas aquí
// 			},
// 		},
// 	}

// 	jsonBody, err := json.Marshal(requestBody)
// 	if err != nil {
// 			return "", err
// 	}
// 	fmt.Println(string(jsonBody)) // Imprime el JSON

// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
// 	if err != nil {
// 			return "", err
// 	}

// 	// 3.  Ya no necesitas establecer el token manualmente. El cliente lo hace automáticamente
// 	// req.Header.Set("Authorization", "Bearer TU_ACCESS_TOKEN")  // Elimina esta línea
// 	req.Header.Set("Content-Type", "application/json")

// 	resp, err := client.Do(req)
// 	if err != nil {
// 			return "", err
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body) // Usa io.ReadAll en lugar de ioutil.ReadAll
// 	if err != nil {
// 			return "", err
// 	}

// 	// ... (Manejo de la respuesta JSON sin cambios)

// 	var response map[string]interface{}
// 	err = json.Unmarshal(body, &response)
// 	if err != nil {
// 			return "", err
// 	}

// 	// Extraer la respuesta del chatbot (con manejo de errores mejorado)
// 	candidates, ok := response["candidates"].([]interface{})
// 	if !ok || len(candidates) == 0 {
// 			return "", fmt.Errorf("no se encontraron candidatos en la respuesta")
// 	}

// 	contentMap, ok := candidates[0].(map[string]interface{})
// 	if !ok {
// 			return "", fmt.Errorf("el primer candidato no es un mapa")
// 	}

// 	content, ok := contentMap["content"]
// 	if !ok {
// 			return "", fmt.Errorf("no se encontró contenido en el candidato")
// 	}

// 	return fmt.Sprintf("%v", content), nil
// }

// ... (Structs VertexAIRequest, Content, Part, Config, SafetySetting sin cambios)

// func GenerateResponse(userMessage string) (map[string]interface{}, error) {
// 	ctx := context.Background()

// 	// 1. Obtener las credenciales de Google Cloud
// 	creds, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/cloud-platform")
// 	if err != nil {
// 		return nil, fmt.Errorf("error obteniendo credenciales: %w", err)
// 	}

// 	// 2. Crear el cliente HTTP con las credenciales
// 	client := oauth2.NewClient(ctx, creds.TokenSource)

// 	url := "https://us-central1-aiplatform.googleapis.com/v1/projects/reelstate-8cc46/locations/us-central1/publishers/google/models/gemini-2.0-flash-exp:generateContent"

// 	requestBody := VertexAIRequest{
// 		Model: "gemini-2.0",
// 		Contents: []Content{
// 			{Role: "user", Parts: []Part{{Text: userMessage}}},
// 		},
// 		Config: Config{
// 			Temperature:        1.0,
// 			TopP:               0.95,
// 			MaxOutputTokens:    8192,
// 			ResponseModalities: []string{"TEXT"}, // Agregado según el nuevo struct
// 		},
// 		SystemInstruction: SystemInstruction{ // Cambio aquí: ahora es un objeto, no un array
// 			Role: "system",
// 			Parts: []Part{
// 				{Text: `Eres un chatbot especializado en bienes raíces. Tu objetivo es ayudar a los usuarios a encontrar propiedades inmobiliarias basadas en criterios específicos como ubicación, precio, tamaño, número de habitaciones y otros detalles.
// 	El chatbot debe realizar las siguientes tareas:
// 	1. **Interpretar la consulta del usuario**: Detecta las palabras clave y los parámetros de la consulta, como la ubicación, el precio, el tamaño y otros detalles relevantes.
// 	2. **Extraer la información**: Basado en la consulta del usuario, debes identificar los parámetros claves que indican las características de la propiedad que está buscando, tales como:
// 	- Ubicación (e.g., "Polanco", "Condesa", "CDMX").
// 	- Precio máximo (e.g., "menos de 2 millones", "máximo 1 millón").
// 	- Área de la propiedad (e.g., "más de 100 metros cuadrados", "más grande que 80 m²").
// 	- Número de habitaciones (e.g., "2 recámaras", "3 habitaciones").
// 	3. **Realizar una búsqueda en la base de datos**: Una vez que los parámetros han sido identificados, devolver esos parámetros en formato JSON, de tal manera que el backend pueda identificar que debe realizar una búsqueda en la base de datos.
// 	4. **Devolver los resultados**: Devuelve los resultados de la búsqueda al usuario en formato natural (por ejemplo, "Encontré 3 propiedades en Polanco por menos de 2 millones de pesos con 3 recámaras y más de 100 m². Aquí están los detalles:").
// 	5. **Formato de la Respuesta**: La respuesta debe incluir información detallada sobre la propiedad, como:
// 	- **Descripción**: Resumen de las características clave de la propiedad.
// 	- **Precio**: El precio de la propiedad.
// 	- **Ubicación**: La ubicación de la propiedad.
// 	- **Área**: El tamaño de la propiedad en metros cuadrados.
// 	- **Número de habitaciones**: El número de habitaciones disponibles.
// 	- **Enlace de contacto**: Un enlace o número de teléfono para contactar con el vendedor o agente.
// 	**Ejemplo de consulta y respuesta**:
// 	- Usuario: "Muéstrame casas en Polanco con 3 recámaras y un precio máximo de 2 millones"
// 	- Chatbot: "He encontrado 2 propiedades en Polanco que coinciden con tu búsqueda. Aquí tienes los detalles:
// 	1. **Casa en Polanco**
// 	 - Precio: 1.8 millones
// 	 - Área: 120 m²
// 	 - Habitaciones: 3 recámaras
// 	 - Descripción: Hermosa casa cerca de parques y tiendas.
	
// 	2. **Departamento en Polanco**
// 	 - Precio: 1.95 millones
// 	 - Área: 110 m²
// 	 - Habitaciones: 3 recámaras
// 	 - Descripción: Departamento moderno, excelente ubicación.
	
// 	¿Te gustaría saber más sobre alguna de estas propiedades?"
	
// 	**Tareas adicionales**: 
// 	- Si la consulta no incluye parámetros claros, puedes pedir al usuario que proporcione más detalles, como el rango de precio, la ubicación y el tipo de propiedad.
// 	- Si no se encuentran propiedades que coincidan con la consulta, el chatbot debe informar al usuario de forma amigable y sugerir otras opciones.
// 	---
// 	Este es un ejemplo de cómo el modelo debería entender y procesar las consultas. El entrenamiento debería incluir múltiples ejemplos de este tipo de interacciones para que el chatbot pueda generalizar y reconocer diferentes variaciones en las consultas.`,
// 				},
// 			},
// 		},
// 		SafetySettings: []SafetySetting{
// 			{Category: "HARM_CATEGORY_HATE_SPEECH", Threshold: "OFF"},
// 			{Category: "HARM_CATEGORY_DANGEROUS_CONTENT", Threshold: "OFF"},
// 			{Category: "HARM_CATEGORY_SEXUALLY_EXPLICIT", Threshold: "OFF"},
// 			{Category: "HARM_CATEGORY_HARASSMENT", Threshold: "OFF"},
// 		},
// 	}

// 	jsonBody, err := json.Marshal(requestBody)
// 	if err != nil {
// 		return nil, err
// 	}

// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
// 	if err != nil {
// 		return nil, err
// 	}

// 	req.Header.Set("Content-Type", "application/json")

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var response map[string]interface{}
// 	err = json.Unmarshal(body, &response)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Extraer la respuesta del chatbot (con manejo de errores mejorado)
// 	candidates, ok := response["candidates"].([]interface{})
// 	if !ok || len(candidates) == 0 {
// 		return nil, fmt.Errorf("no se encontraron candidatos en la respuesta")
// 	}

// 	// contentMap, ok := candidates[0].(map[string]interface{})
// 	// if !ok {
// 	// 	return nil, fmt.Errorf("el primer candidato no es un mapa")
// 	// }

// 	// content, ok := contentMap["content"]
// 	// if !ok {
// 	// 	return nil, fmt.Errorf("no se encontró contenido en el candidato")
// 	// }
// 	// Extraer la respuesta del chatbot en formato JSON
// 	if candidates, ok := response["candidates"].([]interface{}); ok && len(candidates) > 0 {
// 		if contentMap, ok := candidates[0].(map[string]interface{}); ok {
// 			return contentMap, nil
// 		}
// 	}

	
// 	return nil, fmt.Errorf("no se pudo extraer una respuesta válida del modelo")
// }

// package main

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"

// 	"golang.org/x/oauth2"
// 	"golang.org/x/oauth2/google"
// )

// Función para generar respuestas con historial de conversación
// func GenerateResponse(userMessage string, history []Content) (string, error) {
// 	ctx := context.Background()

// 	// Obtener credenciales de Google Cloud
// 	creds, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/cloud-platform")
// 	if err != nil {
// 		return "", fmt.Errorf("error obteniendo credenciales: %w", err)
// 	}

// 	// Crear cliente HTTP con autenticación
// 	client := oauth2.NewClient(ctx, creds.TokenSource)

// 	url := "https://us-central1-aiplatform.googleapis.com/v1/projects/reelstate-8cc46/locations/us-central1/publishers/google/models/gemini-2.0-flash-exp:generateContent"

// 	// Agregar el nuevo mensaje del usuario al historial
// 	history = append(history, Content{Role: "user", Parts: []Part{{Text: userMessage}}})

// 	// Crear cuerpo de la petición
// 	requestBody := VertexAIRequest{
// 		Model:    "gemini-2.0",
// 		Contents: history, // Se envía el historial completo
// 		Config: Config{
// 			Temperature:        1.0,
// 			TopP:               0.95,
// 			MaxOutputTokens:    8192,
// 			ResponseModalities: []string{"TEXT"},
// 		},
// 		SystemInstruction: SystemInstruction{
// 			Role: "system",
// 			Parts: []Part{
// 				{Text: `Eres un chatbot especializado en bienes raíces. Tu objetivo es ayudar a los usuarios a encontrar propiedades inmobiliarias basadas en criterios específicos como ubicación, precio, tamaño, número de habitaciones y otros detalles.
// 	El chatbot debe realizar las siguientes tareas:
// 	1. **Interpretar la consulta del usuario**: Detecta las palabras clave y los parámetros de la consulta, como la ubicación, el precio, el tamaño y otros detalles relevantes.
// 	2. **Extraer la información**: Basado en la consulta del usuario, debes identificar los parámetros claves que indican las características de la propiedad que está buscando, tales como:
// 	- Ubicación (e.g., "Polanco", "Condesa", "CDMX").
// 	- Precio máximo (e.g., "menos de 2 millones", "máximo 1 millón").
// 	- Área de la propiedad (e.g., "más de 100 metros cuadrados", "más grande que 80 m²").
// 	- Número de habitaciones (e.g., "2 recámaras", "3 habitaciones").
// 	3. **Realizar una búsqueda en la base de datos**: Una vez que los parámetros han sido identificados, devolver esos parámetros en formato JSON, de tal manera que el backend pueda identificar que debe realizar una búsqueda en la base de datos.
// 	4. **Devolver los resultados**: Devuelve los resultados de la búsqueda al usuario en formato natural (por ejemplo, "Encontré 3 propiedades en Polanco por menos de 2 millones de pesos con 3 recámaras y más de 100 m². Aquí están los detalles:").
// 	5. **Formato de la Respuesta**: La respuesta debe incluir información detallada sobre la propiedad, como:
// 	- **Descripción**: Resumen de las características clave de la propiedad.
// 	- **Precio**: El precio de la propiedad.
// 	- **Ubicación**: La ubicación de la propiedad.
// 	- **Área**: El tamaño de la propiedad en metros cuadrados.
// 	- **Número de habitaciones**: El número de habitaciones disponibles.
// 	- **Enlace de contacto**: Un enlace o número de teléfono para contactar con el vendedor o agente.
// 	**Ejemplo de consulta y respuesta**:
// 	- Usuario: "Muéstrame casas en Polanco con 3 recámaras y un precio máximo de 2 millones"
// 	- Chatbot: "He encontrado 2 propiedades en Polanco que coinciden con tu búsqueda. Aquí tienes los detalles:
// 	1. **Casa en Polanco**
// 	 - Precio: 1.8 millones
// 	 - Área: 120 m²
// 	 - Habitaciones: 3 recámaras
// 	 - Descripción: Hermosa casa cerca de parques y tiendas.
	
// 	2. **Departamento en Polanco**
// 	 - Precio: 1.95 millones
// 	 - Área: 110 m²
// 	 - Habitaciones: 3 recámaras
// 	 - Descripción: Departamento moderno, excelente ubicación.
	
// 	¿Te gustaría saber más sobre alguna de estas propiedades?"
	
// 	**Tareas adicionales**: 
// 	- Si la consulta no incluye parámetros claros, puedes pedir al usuario que proporcione más detalles, como el rango de precio, la ubicación y el tipo de propiedad.
// 	- Si no se encuentran propiedades que coincidan con la consulta, el chatbot debe informar al usuario de forma amigable y sugerir otras opciones.
// 	---
// 	Este es un ejemplo de cómo el modelo debería entender y procesar las consultas. El entrenamiento debería incluir múltiples ejemplos de este tipo de interacciones para que el chatbot pueda generalizar y reconocer diferentes variaciones en las consultas.`}, // Aquí va tu instrucción completa
// 			},
// 		},
// 		SafetySettings: []SafetySetting{
// 			{Category: "HARM_CATEGORY_HATE_SPEECH", Threshold: "OFF"},
// 			{Category: "HARM_CATEGORY_DANGEROUS_CONTENT", Threshold: "OFF"},
// 			{Category: "HARM_CATEGORY_SEXUALLY_EXPLICIT", Threshold: "OFF"},
// 			{Category: "HARM_CATEGORY_HARASSMENT", Threshold: "OFF"},
// 		},
// 	}

// 	jsonBody, err := json.Marshal(requestBody)
// 	if err != nil {
// 		return "", err
// 	}

// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
// 	if err != nil {
// 		return "", err
// 	}

// 	req.Header.Set("Content-Type", "application/json")

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return "", err
// 	}

// 	var response map[string]interface{}
// 	err = json.Unmarshal(body, &response)
// 	if err != nil {
// 		return "", err
// 	}

// 	// Extraer la respuesta del chatbot
// 	if candidates, ok := response["candidates"].([]interface{}); ok && len(candidates) > 0 {
// 		if contentMap, ok := candidates[0].(map[string]interface{}); ok {
// 			if content, exists := contentMap["content"]; exists {
// 				if contentList, ok := content.([]interface{}); ok && len(contentList) > 0 {
// 					if firstResponse, ok := contentList[0].(map[string]interface{}); ok {
// 						if text, ok := firstResponse["text"].(string); ok {
// 							return text, nil
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}

// 	return "", fmt.Errorf("no se pudo extraer una respuesta válida del modelo")
// }
func GenerateResponse(userMessage string) (string, error) {
	ctx := context.Background()

	// Obtener credenciales de Google Cloud
	creds, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return "", fmt.Errorf("error obteniendo credenciales: %w", err)
	}

	// Crear cliente HTTP con las credenciales
	client := oauth2.NewClient(ctx, creds.TokenSource)

	url := "https://us-central1-aiplatform.googleapis.com/v1/projects/reelstate-8cc46/locations/us-central1/publishers/google/models/gemini-2.0-flash-exp:generateContent"

	requestBody := VertexAIRequest{
		Model: "gemini-2.0",
		Contents: []Content{
			{Role: "user", Parts: []Part{{Text: userMessage}}},
		},
		Config: Config{
			Temperature:        1.0,
			TopP:               0.95,
			MaxOutputTokens:    8192,
			ResponseModalities: []string{"TEXT"},
		},
		SystemInstruction: SystemInstruction{
			Role: "system",
			Parts: []Part{
				{Text: `Eres un chatbot especializado en bienes raíces. Tu objetivo es ayudar a los usuarios a encontrar propiedades inmobiliarias...`}, // Aquí va tu prompt
			},
		},
		SafetySettings: []SafetySetting{
			{Category: "HARM_CATEGORY_HATE_SPEECH", Threshold: "OFF"},
			{Category: "HARM_CATEGORY_DANGEROUS_CONTENT", Threshold: "OFF"},
			{Category: "HARM_CATEGORY_SEXUALLY_EXPLICIT", Threshold: "OFF"},
			{Category: "HARM_CATEGORY_HARASSMENT", Threshold: "OFF"},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	// Extraer la respuesta del chatbot en formato string
	candidates, ok := response["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		return "", fmt.Errorf("no se encontraron candidatos en la respuesta")
	}

	contentMap, ok := candidates[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("el primer candidato no es un mapa válido")
	}

	content, ok := contentMap["content"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("no se encontró contenido en el candidato")
	}

	parts, ok := content["parts"].([]interface{})
	if !ok || len(parts) == 0 {
		return "", fmt.Errorf("no se encontraron partes en la respuesta")
	}

	part, ok := parts[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("el primer elemento en 'parts' no es un mapa válido")
	}

	text, ok := part["text"].(string)
	if !ok {
		return "", fmt.Errorf("no se encontró texto en la respuesta")
	}

	return text, nil
}
