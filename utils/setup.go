package utils

import amqp "github.com/rabbitmq/amqp091-go"

func SetupConsumer(ch *amqp.Channel, topicName string, queueName string) (<-chan amqp.Delivery, error) {
    q, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
    if err != nil {
        return nil, err
    }
    err = ch.QueueBind(q.Name, topicName, "amq.topic", false, nil)
    if err != nil {
        return nil, err
    }
    msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
    return msgs, err
}