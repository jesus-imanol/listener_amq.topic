package consumers

import (
	"encoding/json"
	entities "listernertopic/data"
	"listernertopic/utils"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ProcessHumidityMessages(msgs <-chan amqp.Delivery) {
    for d := range msgs {
        log.Printf("Recibido mensaje de humedad: %s", d.Body)

        var rawData map[string]interface{}
        if err := json.Unmarshal(d.Body, &rawData); err != nil {
            log.Printf("Error al parsear JSON: %s", err)
            continue
        }

        humidity := entities.SensorData{
            ID:       time.Now().Unix(),
            Type:     rawData["type"].(string),
            Quantity: rawData["quantity"].(float64),
            Text:     rawData["text"].(string),
            User:     rawData["user"].(string),
        }

        standardizedJSON, err := json.Marshal(humidity)
        if err != nil {
            log.Printf("Error al crear JSON estandarizado: %s", err)
            continue
        }

        if err := utils.SendToAPI(standardizedJSON); err != nil {
            log.Printf("Error al enviar datos a la API: %s", err)
        } else {
            log.Printf("Datos de humedad enviados exitosamente a la API: %s", standardizedJSON)
        }
    }
}