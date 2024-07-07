package utils

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func LogError(c *gin.Context, err error) {
	_ = c
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	logger.Println(err)
}

func ErrorResponse(c *gin.Context, status int, message Envelope) {

	err := WriteResponse(c, status, message, nil)

	if err != nil {
		LogError(c, err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func NotFoundResponse(c *gin.Context) {
	env := Envelope{"message": "the requested resource could not be found"}
	ErrorResponse(c, http.StatusNotFound, env)
}

func MethodNotAllowedResponse(c *gin.Context) {
	env := Envelope{"message": "the requested method is not allowed"}
	ErrorResponse(c, http.StatusMethodNotAllowed, env)
}

func BadRequestResponse(c *gin.Context, msg string, status int, err error) {
	LogError(c, err)

	env := Envelope{
		"status":     "Bad request",
		"message":    msg,
		"statusCode": status,
	}
	ErrorResponse(c, status, env)
}

func ValidationErrorResponse(c *gin.Context, errors any) {
	// app.logError(c, err)

	env := Envelope{
		"errors": errors,
	}
	ErrorResponse(c, http.StatusUnprocessableEntity, env)
}
