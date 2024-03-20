package main

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

const ADDRESS string = "0.0.0.0"
const PORT string = "8080"

type response struct {
	Message    string `json:"message"`
	HttpStatus int    `json:"http_status"`
	Success    bool   `json:"success"`
}

var pong = response{
	Message:    "pong!",
	HttpStatus: 200,
	Success:    true,
}

func ping(c *gin.Context) {
	c.IndentedJSON(pong.HttpStatus, pong)
}

func main() {
	fullAddress := strings.Join([]string{ADDRESS, PORT}, ":")
	fmt.Print("Starting server on")
	fmt.Println(fullAddress)

	router := gin.Default()
	router.GET("/ping", ping)
	router.Run(fullAddress)
}
