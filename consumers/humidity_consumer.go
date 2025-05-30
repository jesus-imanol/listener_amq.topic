package consumers

import (
	"encoding/json"
	entities "listernertopic/data"
	"listernertopic/utils"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ProcessHumidityMessages(token string, msgs <-chan amqp.Delivery) {
    for d := range msgs {
        log.Printf("Recibido mensaje de humedad: %s", d.Body)

        var rawData map[string]interface{}
        if err := json.Unmarshal(d.Body, &rawData); err != nil {
            log.Printf("Error al parsear JSON: %s", err)
            continue
        }
        nowHour := time.Now()
        mysqlDateTime := nowHour.Format("2006-01-02 15:04:05");
        humidity := entities.SensorData{
            ID:       time.Now().Unix(),
            Type:     rawData["type"].(string),
            Quantity: rawData["quantity"].(float64),
            Text:     rawData["text"].(string),
            User:     rawData["user"].(string),
            CreatedAt: mysqlDateTime,
        }

        standardizedJSON, err := json.Marshal(humidity)
        if err != nil {
            log.Printf("Error al crear JSON estandarizado: %s", err)
            continue
        }
        //err1 := utils.InitEsp32(humidity.User, token)
        //if err1 != nil {
          //  log.Printf("Estado de la inicialización: %s", err)
           // continue
       // }
        if err := utils.SendToAPI(token, standardizedJSON); err != nil {
            log.Printf("Error al enviar datos a la API: %s", err)
        } else {
            log.Printf("Datos de humedad enviados exitosamente a la API: %s", standardizedJSON)
        }
    }
}