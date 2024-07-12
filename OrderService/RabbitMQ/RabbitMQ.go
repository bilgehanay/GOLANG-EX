package rabbitmq

import (
	"fmt"
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

func PublishMessage(body string, messageType string) error {
	err := channel.Publish(
		"",
		"OrderQueue",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
			Type:        messageType,
		})
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}
	fmt.Println("Message published:", body)
	return nil
}

func Close() {
	if err := channel.Close(); err != nil {
		log.Fatalf("Failed to close channel: %v", err)
	}
	if err := conn.Close(); err != nil {
		log.Fatalf("Failed to close connection: %v", err)
	}
}
