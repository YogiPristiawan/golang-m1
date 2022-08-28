package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to create channel")
	defer ch.Close()

	// declare queue as exclusive to make temporary
	// this queue is used to callback queue (replyTo)
	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	failOnError(err, "Failed to declare queue")

	correlationId := randomString(5)
	fibonacciIdx := getFibonacciIndex(os.Args)

	err = ch.Publish(
		"",
		"rpc_queue",
		false,
		false,
		amqp.Publishing{
			Body:          []byte(strconv.Itoa(fibonacciIdx)),
			CorrelationId: correlationId,
			ReplyTo:       q.Name,
			ContentType:   "text/plain",
		},
	)
	failOnError(err, "Failed to publish into queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to consume queue")

	for msg := range msgs {
		if correlationId == msg.CorrelationId {
			res, err := strconv.Atoi(string(msg.Body))
			failOnError(err, "Failed to print response")

			fmt.Printf("[v] Response fibonacci ke- %d : %d\n", fibonacciIdx, res)
			break
		}
	}

}

func randomString(length int) (str string) {
	for i := 0; i < length; i++ {
		str += string(randInt(65, 90))
	}
	return
}

func randInt(min int, max int) byte {
	return byte(min + (rand.Intn(max - min)))
}

func getFibonacciIndex(args []string) int {
	if len(args) < 2 || args[1] == "" {
		log.Printf("Index of fibonacci required")
		os.Exit(1)
	}

	num, err := strconv.Atoi(args[1])
	failOnError(err, "Parameter must be a number")
	return num
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", err, msg)
	}
}
