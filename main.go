package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// SensorData estructura para los datos del sensor
type SensorData struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
	Text     string `json:"text"`
	Message  string `json:"message"`
	Humedad  int    `json:"humedad"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func sendToAPI(data []byte) error {
	// Reemplaza esta URL con la de tu API
	apiURL := "http://127.0.0.1:4000/v1/message/"
	
	// Crear la solicitud HTTP
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	
	// Configurar encabezados
	req.Header.Set("Content-Type", "application/json")
	
	// Crear cliente HTTP con timeout
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	
	// Enviar solicitud
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	// Registrar respuesta
	log.Printf("API response status: %s", resp.Status)
	return nil
}

func main() {
	conn, err := amqp.Dial("amqp://Somer:140823schmesom2018@44.209.159.224:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	
	err = ch.ExchangeDeclare(
		"amq.topic", // name
		"topic",     // type
		true,        // durable
		false,       // auto-deleted
		false,       // internal
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare an exchange")
	
	q, err := ch.QueueDeclare(
		"humedad", // nombre
		true,      // durable
		false,     // eliminar cuando no se use
		false,     // exclusivo
		false,     // no esperar
		nil,       // argumentos
	)
	failOnError(err, "Failed to declare a queue")
	
	// Si no se proporcionan claves de enrutamiento, usar el topic por defecto
	bindingKeys := os.Args[1:]
	if len(bindingKeys) == 0 {
		bindingKeys = []string{"systemHumidity.mqtt"} // Topic por defecto
	}
	
	for _, s := range bindingKeys {
		log.Printf("Binding queue %s to exchange %s with routing key %s",
			q.Name, "amq.topic", s)
		err = ch.QueueBind(
			q.Name,      // queue name
			s,           // routing key
			"amq.topic", // exchange
			false,
			nil)
		failOnError(err, "Failed to bind a queue")
	}
	
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")
	
	var forever chan struct{}
	
	go func() {
		for d := range msgs {
			log.Printf("Received message: %s", d.Body)
			
			// Parsear el mensaje JSON
			var sensorData SensorData
			if err := json.Unmarshal(d.Body, &sensorData); err != nil {
				log.Printf("Error parsing JSON: %s", err)
				continue
			}
			
			// Registrar los datos del sensor
			log.Printf("Sensor data - Type: %s, Quantity: %d, Text: %s",
				sensorData.Type, sensorData.Quantity, sensorData.Text)
			
			// Enviar datos a la API
			if err := sendToAPI(d.Body); err != nil {
				log.Printf("Error sending data to API: %s", err)
			} else {
				log.Printf("Data successfully sent to API")
			}
		}
	}()
	
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	}