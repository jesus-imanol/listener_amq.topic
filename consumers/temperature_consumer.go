package consumers

import (
	"encoding/json"
	entities "listernertopic/data"
	"listernertopic/utils"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ProcessTemperatureMessages(msgs <-chan amqp.Delivery) {
    for d := range msgs {
        log.Printf("Recibido mensaje de temperatura: %s", d.Body)

        var rawData map[string]interface{}
        if err := json.Unmarshal(d.Body, &rawData); err != nil {
            log.Printf("Error al parsear JSON: %s", err)
            continue
        }
        nowHour := time.Now()
        mysqlDateTime := nowHour.Format("2006-01-02 15:04:05");
        log.Printf("Formato de fecha y hora: %s", mysqlDateTime)
        temperature := entities.SensorData{
            ID:       time.Now().Unix(),
            Type:     rawData["type"].(string),
            Quantity: rawData["quantity"].(float64),
            Text:     rawData["text"].(string),
            User:     rawData["user"].(string),
            CreatedAt: mysqlDateTime,
        }

        standardizedJSON, err := json.Marshal(temperature)
        if err != nil {
            log.Printf("Error al crear JSON estandarizado: %s", err)
            continue
        }

        if err := utils.SendToAPI(standardizedJSON); err != nil {
            log.Printf("Error al enviar datos a la API: %s", err)
        } else {
            log.Printf("Datos de temperatura enviados exitosamente a la API: %s", standardizedJSON)
        }
    }
}