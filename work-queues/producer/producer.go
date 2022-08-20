package main

import (
	"log"
	"os"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to create channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"tasks",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare queue")

	// tage arguments from command line
	msg := strings.Join(os.Args[1:], " ")

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(msg),
		},
	)
	failOnError(err, "Failed to publish message")

	log.Printf("[x] Sent %s", msg)
}
