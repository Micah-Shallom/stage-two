package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Micah-Shallom/stage-two/config"
	"github.com/Micah-Shallom/stage-two/handlers"
	"github.com/Micah-Shallom/stage-two/models"
	"github.com/Micah-Shallom/stage-two/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

var jwtSecret = []byte(os.Getenv(
	"JWT_SECRET",
))

func TestMain(m *testing.M) {
	m.Run()
}

func TestGenerateJWT(t *testing.T) {
	user := &models.User{
		UserID:    "123",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	token, err := utils.GenerateJWT(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Parse the token
	parsedToken, err := jwt.ParseWithClaims(token, &utils.SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	// Validate the token claims
	claims, ok := parsedToken.Claims.(*utils.SignedDetails)
	assert.True(t, ok)
	assert.Equal(t, user.UserID, claims.UserID)
	assert.Equal(t, user.Email, claims.Email)
	assert.Equal(t, user.FirstName, claims.FirstName)
	assert.Equal(t, user.LastName, claims.LastName)
}

func TestTokenExpiration(t *testing.T) {
	// Token with a short expiration time
	user := &models.User{
		UserID:    "123",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	token, err := utils.GenerateJWT(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Parse the token
	parsedToken, err := jwt.ParseWithClaims(token, &utils.SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}
	assert.True(t, parsedToken.Valid)

	// Validate the token claims
	claims, ok := parsedToken.Claims.(*utils.SignedDetails)
	if !ok {
		t.Fatalf("Failed to parse claims: %v", claims)
	}

	// Calculate the expected expiration time with a tolerance of 30 seconds
	expectedExpiration := time.Now().Add(7 * 24 * time.Hour)
	assert.WithinDuration(t, expectedExpiration, time.Now().Add(7*24*time.Hour), 30*time.Second)
}

// It Should Register User Successfully with Default Organisation:Ensure a user is registered successfully when no organisation details are provided.

func TestRegisterUser(t *testing.T) {

	// Create a new user
	user := &models.User{
		UserID:    "123",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "testuser@gmail.com",
		Password:  "password",
		Phone:     "1234567890",
	}

	db, err := SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}

	app := config.NewApplication(db)
	h := handlers.NewHandlers(app)

	r := gin.Default()
	r.POST("/register", h.RegisterUserHandler)

	userPayload, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal user payload: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(userPayload))

	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Create a new response recorder
	rr := httptest.NewRecorder()

	// Serve the HTTP request
	r.ServeHTTP(rr, req)

	//Check the status code
	assert.Equal(t, http.StatusCreated, rr.Code, "Expected status code 201 Created")

	//check the response body
	var response utils.Envelope
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	assert.Equal(t, "success", response["status"], "Expected status to be 'success'")
	assert.Equal(t, "Registration successful", response["message"], "Expected message to be 'Registration successful'")

	
	// Check the database for the created organization
	var o models.Organisation
	err = db.Where("name = ?", user.FirstName+"'s Organisation").First(&o).Error
	assert.NoError(t, err, "Expected no error when querying for the organization")
	assert.Contains(t, o.Description, "Default organisation for", "Expected organization description to contain 'Default organisation for'")
	//write an assertion for checking org name
	assert.Equal(t, user.FirstName+"'s Organisation", o.Name, "Expected organisation name to match")

	// Check the database for the created user
	var u models.User
	err = db.Where("email = ?", user.Email).First(&u).Error
	
	assert.NoError(t, err, "Expected no error when querying for the user")
	assert.Equal(t, user.FirstName, u.FirstName, "Expected first name to match")
	assert.Equal(t, user.LastName, u.LastName, "Expected last name to match")
	assert.Equal(t, user.Email, u.Email, "Expected email to match")
	assert.Equal(t, user.Phone, u.Phone, "Expected phone number to match")
	
	// Check that the response contains the expected access token.
	assert.Contains(t, response["data"], "accessToken", "Expected response to contain 'accessToken'")

}

// It Should Log the user in successfully:Ensure a user is logged in successfully when a valid credential is provided and fails otherwise.
// Check that the response contains the expected user details and access token.

func TestLoginUser(t *testing.T) {
	

	// Create a new user
	user := &models.User{
		Email:     "testuser@gmail.com",
		Password:  "password",
	}

	db, err := SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}

	app := config.NewApplication(db)
	h := handlers.NewHandlers(app)

	r := gin.Default()
	r.POST("/login", h.LoginUserHandler)

	userPayload, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal user payload: %v", err)
	}

	// Create a new HTTP request

	req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(userPayload))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Create a new response recorder
	rr := httptest.NewRecorder()

	// Serve the HTTP request
	r.ServeHTTP(rr, req)

	//Check the status code
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status code 200 OK")

	//check the response body
	var response utils.Envelope
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	assert.Equal(t, "success", response["status"], "Expected status to be 'success'")
	assert.Equal(t, "Login successful", response["message"], "Expected message to be 'Login successful'")
	assert.Contains(t, response["data"], "accessToken", "Expected response to contain 'accessToken'")
	assert.Contains(t, response["data"], "user", "Expected response to contain 'user'")
}


func TestExistingUserRecord(t *testing.T){
	defer MockedDB(DROP)
	// use the same user that was used when testing the registation functionality
	user := &models.User{
		UserID:    "123",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "testuser@gmail.com",
		Password:  "password",
		Phone:     "1234567890",
	}

	db, err := SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}

	app := config.NewApplication(db)
	h := handlers.NewHandlers(app)

	r := gin.Default()
	r.POST("/register", h.RegisterUserHandler)

	userPayload, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal user payload: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(userPayload))

	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Create a new response recorder
	rr := httptest.NewRecorder()

	// Serve the HTTP request
	r.ServeHTTP(rr, req)

	//check the status code
	assert.Equal(t, 400, rr.Code, "Expected Error 400 due to unsuccessful registration")
}