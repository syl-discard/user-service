package middlewares

import (
	messenger "discard/user-service/pkg/messenger"
	"discard/user-service/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	amqp "github.com/rabbitmq/amqp091-go"
)

func ForwardUserDeletionRequestMiddleware(channel *amqp.Channel, queueName string, message string) gin.HandlerFunc {
	return func(context *gin.Context) {
		var user models.User
		if err := context.ShouldBindBodyWith(&user, binding.JSON); err != nil {
			context.AbortWithStatusJSON(
				http.StatusBadRequest, models.Response{
					Message:    "Failed to send deletion request: " + err.Error(),
					HttpStatus: http.StatusBadRequest,
					Success:    false,
				})
			return
		}

		messenger.Message(channel, queueName, []byte(message+user.ID))

		context.Next()
	}
}
