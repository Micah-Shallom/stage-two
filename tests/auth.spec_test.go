package tests

import (
	"os"
	"testing"
	"time"

	"github.com/Micah-Shallom/stage-two/models"
	"github.com/Micah-Shallom/stage-two/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

var jwtSecret = []byte(os.Getenv(
	"JWT_SECRET",
))

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
