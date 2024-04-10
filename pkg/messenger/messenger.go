package messenger

import (
	"context"
	"time"

	logger "discard/user-service/pkg/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Message(channel *amqp.Channel, queueName string, message []byte) {
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
			Body:        message,
		}, // message
	)

	logger.FailOnError(err, "Failed to send a message")
}
