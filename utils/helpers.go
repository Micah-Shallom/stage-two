package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Envelope map[string]any

type validationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func userFriendlyMessage(tag string) string {
	messages := map[string]string{
		"required": "is required",
		"email":    "must be a valid email address",
		"alpha":    "should only contain alphabetic characters",
		"alphanum": "should only contain alphanumeric characters",
		"min":      "should be at least 8 characters long",
		"max":      "should be at most 15 characters long",
		"eq":       "should be equal to %d",
		"ne":       "should not be equal to %d",
		"lt":       "should be less than %d",
		"lte":      "should be less than or equal to %d",
		"gt":       "should be greater than %d",
		"gte":      "should be greater than or equal to %d",
		"len":      "should be %d characters long",
	}
	return messages[tag]
}

func ReadRequest(c *gin.Context, dst interface{}) (bool, []validationError) {
	var validationErrors []validationError

	// Decode the request body into the destination struct
	err := c.ShouldBindJSON(dst)
	if err != nil {
		// Check if it's a JSON unmarshal type error
		var jsonErr *json.UnmarshalTypeError
		if errors.As(err, &jsonErr) {
			_, errors := validateFields(dst)
			validationErrors = append(validationErrors, errors...)
		} else {
			// Generic JSON parsing error
			errorDetail := validationError{
				Field:   "body",
				Message: "Invalid JSON format",
			}
			validationErrors = append(validationErrors, errorDetail)
		}
		return false, validationErrors
	}

	// Validate the request payload
	return validateFields(dst)
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

func validateFields(user interface{}) (bool, []validationError) {
	var validationErrors []validationError

	err := validate.Struct(user)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return false, validationErrors // Invalid validation error
		}

		for _, err := range err.(validator.ValidationErrors) {
			errorDetail := validationError{
				Field: err.Field(),
				// Message: fmt.Sprintf("Validation failed on '%s' for condition '%s'", err.Field(), err.Tag()),
				Message: fmt.Sprintf("The %s field %s", err.Field(), userFriendlyMessage(err.Tag())),
			}
			validationErrors = append(validationErrors, errorDetail)
		}

		return false, validationErrors
	}

	return true, validationErrors
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
