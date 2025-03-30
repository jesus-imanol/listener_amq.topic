package main

import (
	"listernertopic/consumers"
	"listernertopic/utils"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
    conn, err := amqp.Dial("amqp://Somer:140823schmesom2018@44.209.159.224:5672/")
    if err != nil {
        log.Fatalf("Error al conectar a RabbitMQ: %s", err)
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("Error al abrir canal RabbitMQ: %s", err)
    }
    defer ch.Close()

    humidityMsgs, err := utils.SetupConsumer(ch, "systemHumidity.mqtt", "humedad")
    if err != nil {
        log.Fatalf("Error al configurar consumidor de humedad: %s", err)
    }

    temperatureMsgs, err := utils.SetupConsumer(ch, "systemTemperature.mqtt", "temperatura")
    if err != nil {
        log.Fatalf("Error al configurar consumidor de temperatura: %s", err)
    }
    lightMsgs, err := utils.SetupConsumer(ch, "ambientLightSensor.mqtt", "luz")
    if err != nil {
        log.Fatalf("Error al configurar consumidor de luz: %s", err)
    }
   //se ejecutan de manera paralela
    go consumers.ProcessHumidityMessages(humidityMsgs)
    go consumers.ProcessTemperatureMessages(temperatureMsgs)
    go consumers.ProcessLightMessages(lightMsgs)
    
    log.Println("Esperando mensajes. Presiona CTRL+C para salir.")
    var forever chan struct{}
    <-forever
}