// 
package services

import (
    "encoding/json"
    "fmt"

    "github.com/streadway/amqp"
)

// Función para enviar mensaje a RabbitMQ
func SendVideoProcessingTask(typeVideo int, videoPath string, requiresAudio bool, audioPath string) error {
    // Conectar a RabbitMQ
    conn, err := amqp.Dial("amqp://ReelState:ReelState2024@rabbitmq:5672/myvhost")
    if err != nil {
        return fmt.Errorf("error al conectar con RabbitMQ: %v", err)
    }
    defer conn.Close()

    // Crear un canal a partir de la conexión
    ch, err := conn.Channel()
    if err != nil {
        return fmt.Errorf("error al abrir el canal en RabbitMQ: %v", err)
    }
    defer ch.Close()

    // Declarar la cola si aún no existe
    queue, err := ch.QueueDeclare(
        "celery", // Nombre de la cola
        true,                 // Duradera
        false,                // No autoeliminar
        false,                // Exclusiva
        false,                // No esperar
        nil,                  // Argumentos
    )
    if err != nil {
        return fmt.Errorf("error al declarar la cola en RabbitMQ: %v", err)
    }

    // Definir el contenido de la tarea
    task := map[string]interface{}{
        "task": "api.task.process_video_task",
        "args": []interface{}{videoPath, requiresAudio, audioPath},
    }

    // Convertir el contenido de la tarea a formato JSON
    taskBytes, err := json.Marshal(task)
    if err != nil {
        return fmt.Errorf("error al serializar la tarea en JSON: %v", err)
    }

    // Publicar la tarea en RabbitMQ
    err = ch.Publish(
        "",           // Exchange (vacío para usar la cola por defecto)
        queue.Name,   // Routing key (nombre de la cola)
        false,        // Mandatorio
        false,        // Inmediato
        amqp.Publishing{
            ContentType: "application/json",
            Body:        taskBytes,
        },
    )
    if err != nil {
        return fmt.Errorf("error al publicar la tarea en RabbitMQ: %v", err)
    }

    fmt.Println("Tarea enviada a RabbitMQ:", string(taskBytes))
    return nil
}
