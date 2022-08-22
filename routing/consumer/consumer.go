package main

import (
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to create channel")
	defer ch.Close()

	// declare queues as exclusive to make temporary
	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	failOnError(err, "Failed to declare queue")

	// declare exchange
	err = ch.ExchangeDeclare(
		"logs_direct",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to create exchange")

	// bind exchange into queue
	err = ch.QueueBind(
		q.Name,
		getSeverity(os.Args),
		"logs_direct",
		false,
		nil,
	)
	failOnError(err, "Failed to bind exchange")

	//  consume
	msg, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to consume message")

	go func() {
		for d := range msg {
			log.Printf("[x] Received message: %s", d.Body)
		}
	}()

	var forever chan struct{}

	<-forever
}

func getSeverity(args []string) string {
	if len(args) < 2 || args[1] == "" {
		return "info"
	}

	if args[1] == "info" || args[1] == "warning" || args[1] == "critical" {
		return args[1]
	}
	fmt.Println("allowed severity is info, warning, or critical")
	os.Exit(1)
	return ""
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", err, msg)
	}
}
