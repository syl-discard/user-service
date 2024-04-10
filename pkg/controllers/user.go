package controllers

import (
	"discard/user-service/pkg/models"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func DeleteUser(context *gin.Context) {
	var user models.User
	if err := context.ShouldBindBodyWith(&user, binding.JSON); err != nil {
		context.AbortWithStatusJSON(
			http.StatusBadRequest, models.Response{
				Message:    "Failed to bind JSON: " + err.Error(),
				HttpStatus: http.StatusBadRequest,
				Success:    false,
			})
		return
	}

	context.IndentedJSON(http.StatusOK, models.Response{
		Message:    "Successfully sent request to delete user: " + user.ID,
		HttpStatus: http.StatusOK,
		Success:    true,
	})
}
