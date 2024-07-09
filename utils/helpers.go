package utils

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

type Envelope map[string]any

// type validationError struct {
// 	Field   string `json:"field"`
// 	Message string `json:"message"`
// }

func ReadRequest(c *gin.Context, dst interface{}) error {

	// Decode the request body into the destination struct
	err := c.ShouldBindJSON(dst)
	if err != nil {
		BadRequestResponse(c, "Registration unsuccessful", 400, err)
	}
	return err
}

func WriteResponse(c *gin.Context, status int, data Envelope, headers http.Header) error {
	// Marshal the data to JSON with indentation
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// Add a newline for readability (optional)
	js = append(js, '\n')

	// Loop through the map of headers and set them using Gin's context
	for key, value := range headers {
		for _, v := range value {
			c.Header(key, v)
		}
	}

	// Set the Content-Type header explicitly, though Gin does this automatically for JSON responses
	c.Header("Content-Type", "Application/json")

	// Use Gin's context to send the status code and JSON response
	c.Data(status, "Application/json", js)

	return nil
}

func GenerateUUID() string {
	UserID := uuid.New()
	return UserID.String()
}

func HashPassword(password string) (string, error) {
	// Hash the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hashedPassword), err
}

func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	LogError(nil, err)
	return err == nil
}

func IsAuthenticated(userID string, authorizedID any) bool {
	return userID == authorizedID
}
