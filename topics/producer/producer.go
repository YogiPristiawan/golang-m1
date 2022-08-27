package main

import (
	"log"
	"os"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to create channel")
	defer ch.Close()

	// declare exchange
	err = ch.ExchangeDeclare(
		"logs_topic",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to create exchange")

	// publish to exchange
	err = ch.Publish(
		"logs_topic",
		getTopicKey(os.Args),
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(getMessage(os.Args)),
		},
	)
	failOnError(err, "Failed to publish message")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", err, msg)
	}
}

func getTopicKey(args []string) string {
	if len(args) < 3 || args[2] == "" {
		return "sentry.info"
	}

	// validate input
	if args[1] != "sentry" && args[1] != "slack" && args[1] != "*" && args[1] != "#" {
		log.Println("invalid log driver")
		os.Exit(1)
	}

	if args[2] != "info" && args[2] != "warning" && args[2] != "critical" && args[2] != "*" && args[2] != "#" {
		log.Println("invalid log type")
		os.Exit(1)
	}

	return strings.Join([]string{args[1], args[2]}, ".")
}

func getMessage(args []string) string {
	if len(args) < 4 || args[3] == "" {
		return "default: Hello world"
	}

	return strings.Join(args[1:], " ")
}
