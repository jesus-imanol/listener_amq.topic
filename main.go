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

type SensorData struct {
		Type     string      `json:"type"`
		Quantity interface{} `json:"quantity"`
		Text     string      `json:"text"`
	}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func sendToAPI(data []byte) error {
	apiURL := "http://127.0.0.1:4000/v1/message/"
	
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	log.Printf("API response status: %s", resp.Status)
	return nil
}

func setupConsumer(ch *amqp.Channel, topicName string, queueName string) (<-chan amqp.Delivery, error) {
	// Declarar la cola
	q, err := ch.QueueDeclare(
		queueName, // nombre de la cola
		true,      // durable
		false,     // eliminar cuando no se use
		false,     // exclusivo
		false,     // no esperar
		nil,       // argumentos
	)
	if err != nil {
		return nil, err
	}
	
	// Enlazar la cola al exchange con la routing key
	err = ch.QueueBind(
		q.Name,      // queue name
		topicName,   // routing key
		"amq.topic", // exchange
		false,
		nil)
	if err != nil {
		return nil, err
	}
	
	log.Printf("Queue %s bound to exchange amq.topic with routing key %s", q.Name, topicName)
	
	// Configurar el consumidor
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	if err != nil {
		return nil, err
	}
	
	return msgs, nil
}

// processMessages procesa los mensajes recibidos
func processMessages(msgs <-chan amqp.Delivery, sensorType string) {
	for d := range msgs {
		log.Printf("Received %s message: %s", sensorType, d.Body)
		
		// Parsear el mensaje JSON recibido
		var rawData map[string]interface{}
		if err := json.Unmarshal(d.Body, &rawData); err != nil {
			log.Printf("Error parsing JSON: %s", err)
			continue
		}
		
		// Crear la estructura estandarizada
		sensorData := SensorData{
			Type:     rawData["type"].(string),
			Quantity: rawData["quantity"],
			Text:     rawData["text"].(string),
		}
		
		// Convertir de nuevo a JSON para asegurar el formato correcto
		standardizedJSON, err := json.Marshal(sensorData)
		if err != nil {
			log.Printf("Error creating standardized JSON: %s", err)
			continue
		}
		
		// Registrar los datos procesados
		log.Printf("%s data (standardized) - Type: %s, Quantity: %v, Text: %s",
			sensorType, sensorData.Type, sensorData.Quantity, sensorData.Text)
		
		// Enviar datos estandarizados a la API
		if err := sendToAPI(standardizedJSON); err != nil {
			log.Printf("Error sending %s data to API: %s", sensorType, err)
		} else {
			log.Printf("%s data successfully sent to API", sensorType)
		}
	}
}

func main() {
	// Conectar a RabbitMQ
	conn, err := amqp.Dial("amqp://Somer:140823schmesom2018@44.209.159.224:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	
	// Abrir canal
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	
	// Declarar exchange
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
	
	// Definir topics predeterminados o usar los proporcionados como argumentos
	humidityTopic := "systemHumidity.mqtt"
	temperatureTopic := "systemTemperature.mqtt"
	
	// Si se proporcionan argumentos, usar el primero para humedad y el segundo para temperatura
	if len(os.Args) > 1 {
		humidityTopic = os.Args[1]
	}
	if len(os.Args) > 2 {
		temperatureTopic = os.Args[2]
	}
	
	// Configurar consumidor para humedad
	humidityMsgs, err := setupConsumer(ch, humidityTopic, "humedad")
	failOnError(err, "Failed to setup humidity consumer")
	
	// Configurar consumidor para temperatura
	temperatureMsgs, err := setupConsumer(ch, temperatureTopic, "temperatura")
	failOnError(err, "Failed to setup temperature consumer")
	
	// Procesar mensajes de humedad en una goroutine
	go processMessages(humidityMsgs, "Humidity")
	
	// Procesar mensajes de temperatura en una goroutine
	go processMessages(temperatureMsgs, "Temperature")
	
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	
	// Canal para mantener la aplicación corriendo
	var forever chan struct{}
	<-forever
}