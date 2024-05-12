package main

import (
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	controllers "discard/user-service/pkg/controllers"
	logger "discard/user-service/pkg/logger"
	messenger "discard/user-service/pkg/messenger"
	middlewares "discard/user-service/pkg/middlewares"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	godotenv.Load()
	var (
		ADDRESS                 string = "0.0.0.0"
		PORT                    string = "8080"
		RABBITMQ_SERVER_ADDRESS string = os.Getenv("RABBITMQ_SERVER_ADDRESS")
	)

	connected := false
	var activeConnection *amqp.Connection = nil
	for !connected {
		conn, err := amqp.Dial(RABBITMQ_SERVER_ADDRESS)
		if err != nil {
			logger.WARN.Println("Failed to connect to RabbitMQ, retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
		} else {
			activeConnection = conn
			connected = true
		}
	}
	defer activeConnection.Close()
	logger.LOG.Println("Successfully connected to RabbitMQ!")

	channel, err := activeConnection.Channel()
	logger.FailOnError(err, "Failed to create a channel")
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"delete-user", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	logger.FailOnError(err, "Failed to declare a queue")

	logger.LOG.Printf("Successfully declared a queue: %v\n", queue)

	message := []byte("Hello, World!")
	messenger.Message(channel, queue.Name, message)
	logger.LOG.Printf("Successfully published a message to RabbitMQ: %v\n", string(message[:]))

	fullAddress := strings.Join([]string{ADDRESS, PORT}, ":")
	logger.LOG.Printf("Starting API server on %v...\n", fullAddress)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET(
		"/ping",
		controllers.Ping,
	)
	router.POST(
		"/delete-user",
		middlewares.ForwardUserDeletionRequestMiddleware(channel, queue.Name, "Deletion request for user: "),
		controllers.DeleteUser,
	)

	logger.LOG.Printf("API server started on %v!\n", fullAddress)
	logger.FailOnError(router.Run(fullAddress), "Failed to run the server")
}
