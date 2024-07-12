package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

var conn *amqp.Connection
var channel *amqp.Channel

func InitRabbitMQ() {
	var err error
	conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	channel, err = conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	_, err = channel.QueueDeclare(
		"OrderQueue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}
}

func ConsumeMessages(handler func(amqp.Delivery)) {
	msgs, err := channel.Consume(
		"OrderQueue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	go func() {
		for d := range msgs {
			handler(d)
		}
	}()
}

func Close() {
	if err := channel.Close(); err != nil {
		log.Fatalf("Failed to close channel: %v", err)
	}
	if err := conn.Close(); err != nil {
		log.Fatalf("Failed to close connection: %v", err)
	}
}
