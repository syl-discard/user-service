package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"

	controllers "discard/user-service/pkg/controllers"
)

const (
	ADDRESS                 string = "0.0.0.0"
	PORT                    string = "8080"
	RABBITMQ_SERVER_ADDRESS string = "amqp://guest:guest@rabbitmq:5672/"
)

var (
	WARN  = log.New(os.Stderr, "[WARNING]\t", log.LstdFlags|log.Lmsgprefix)
	ERROR = log.New(os.Stderr, "[ERROR]\t", log.LstdFlags|log.Lmsgprefix)
	LOG   = log.New(os.Stdout, "[INFO]\t", log.LstdFlags|log.Lmsgprefix)
)

func failOnError(err error, msg string) {
	if err != nil {
		ERROR.Printf("%s: %s\n", msg, err)
		panic(err)
	}
}

func message(channel *amqp.Channel, queueName string, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := channel.PublishWithContext(
		ctx,       // context
		"",        // exchange
		queueName, // key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		}, // message
	)
	failOnError(err, "Failed to publish a message")
}

func main() {
	connected := false
	var activeConnection *amqp.Connection = nil
	for !connected {
		conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err != nil {
			WARN.Println("Failed to connect to RabbitMQ, retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
		} else {
			activeConnection = conn
			connected = true
		}
	}
	defer activeConnection.Close()
	if !connected {
		ERROR.Println("Failed to connect to RabbitMQ after multiple retries, exiting...")
		os.Exit(1)
	}

	LOG.Println("Successfully connected to RabbitMQ!")

	ch, err := activeConnection.Channel()
	failOnError(err, "Failed to create a channel")
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	LOG.Printf("Successfully declared a queue: %v\n", queue)

	msg := "Hello, World!"
	message(ch, queue.Name, msg)
	LOG.Printf("Successfully published a message to RabbitMQ: %v\n", msg)

	fullAddress := strings.Join([]string{ADDRESS, PORT}, ":")
	LOG.Printf("Starting API server on %v...\n", fullAddress)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/ping", controllers.Ping)

	LOG.Printf("API server started on %v!\n", fullAddress)
	failOnError(router.Run(fullAddress), "Failed to run the server")
}
