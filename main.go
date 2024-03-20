package main

import (
	"context"
	"discard/user-service/pkg/models"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
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

var pong = models.Response{
	Message:    "pong!",
	HttpStatus: 200,
	Success:    true,
}

func ping(c *gin.Context) {
	c.IndentedJSON(pong.HttpStatus, pong)
}

func failOnError(err error, msg string) {
	if err != nil {
		ERROR.Printf("%s: %s\n", msg, err)
		panic(err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	LOG.Println("Successfully connected to RabbitMQ!")

	ch, err := conn.Channel()
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello World!"
	err = ch.PublishWithContext(
		ctx,     // context
		"",      // exchange
		"hello", // key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		}, // message
	)
	failOnError(err, "Failed to publish a message")

	LOG.Printf("Successfully published a message to RabbitMQ: %v\n", body)

	fullAddress := strings.Join([]string{ADDRESS, PORT}, ":")
	LOG.Printf("Starting API server on %v...\n", fullAddress)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/ping", ping)

	LOG.Printf("API server started on %v!\n", fullAddress)
	failOnError(router.Run(fullAddress), "Failed to run the server")
}
