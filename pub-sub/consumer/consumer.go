package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", err, msg)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to create channel")
	defer ch.Close()

	// declare queue as exclusive to make temporary
	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // auto delete
		true,  // exclusive
		false, // no wait
		nil,   // args
	)
	failOnError(err, "Failed to declare queue")

	err = ch.ExchangeDeclare(
		"logs",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare exchange")

	// bind queue into specific exchanges
	err = ch.QueueBind(
		q.Name,
		"",
		"logs",
		false,
		nil,
	)
	failOnError(err, "Failed to create queue")

	msg, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to consume")

	var forever <-chan struct{}

	go func() {
		for d := range msg {
			log.Printf("[x] Received: %s", d.Body)
		}
	}()

	<-forever
}
