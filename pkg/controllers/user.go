package controllers

import "github.com/gin-gonic/gin"

func DeleteUser(context *gin.Context) {
	context.IndentedJSON(pong.HttpStatus, pong)
}
