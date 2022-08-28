package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var clear map[string]func()

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to create connection")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to create channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"rpc_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to create queue")

	err = ch.Qos(1, 0, false)
	failOnError(err, "Failed to crate Qos")

	msg, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	var forever chan struct{}

	go func() {
		for d := range msg {
			number, err := strconv.Atoi(string(d.Body))
			failOnError(err, "Body is invalid")
			fibonacci := strconv.Itoa(fib(number))

			// genrate log message
			delayedInSeconds := rand.Intn(10)
			for i := delayedInSeconds; i >= 0; i-- {
				logMsg := fmt.Sprintf("[.] Processing fn => fib(%d) : estimated %ds", number, i)
				printLog(logMsg)

				if i > 0 {
					time.Sleep(1 * time.Second)
				}
			}

			err = ch.Publish(
				"",
				d.ReplyTo,
				false,
				false,
				amqp.Publishing{
					Body:          []byte(fibonacci),
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
				},
			)
			failOnError(err, "Failed to publish to ReplyTo queue")

			d.Ack(false)
		}
	}()

	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", err, msg)
	}
}

func fib(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fib(n-1) + fib(n-2)
	}
}

func printLog(str string) {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
		fmt.Print(str)
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}
