package models

// for binding validation, see: https://github.com/go-playground/validator

type Response struct {
	Message    string `json:"message"`
	HttpStatus int    `json:"http_status"`
	Success    bool   `json:"success"`
}

type User struct {
	ID string `json:"id" binding:"required,uuid"`
}

type Message []byte
