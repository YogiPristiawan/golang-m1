package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQs")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to create channel")
	defer ch.Close()

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
	failOnError(err, "Failed to declare exchange")

	// publish into exchange
	err = ch.Publish(
		"logs_direct",
		getSeverity(os.Args),
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(getMessage(os.Args)),
		},
	)
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

func getMessage(args []string) string {
	if len(args) < 3 || args[2] == "" {
		return "Hello world!"
	}

	return strings.Join(args[2:], " ")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", err, msg)
	}
}
