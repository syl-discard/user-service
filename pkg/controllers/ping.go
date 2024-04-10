package controllers

import (
	"discard/user-service/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

var pong = models.Response{
	Message:    "pong!",
	HttpStatus: http.StatusOK,
	Success:    true,
}

func Ping(context *gin.Context) {
	context.IndentedJSON(pong.HttpStatus, pong)
}
