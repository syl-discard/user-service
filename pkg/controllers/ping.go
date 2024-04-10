package controllers

import (
	"discard/user-service/pkg/models"

	"github.com/gin-gonic/gin"
)

var pong = models.Response{
	Message:    "pong!",
	HttpStatus: 200,
	Success:    true,
}

func Ping(context *gin.Context) {
	context.IndentedJSON(pong.HttpStatus, pong)
}
