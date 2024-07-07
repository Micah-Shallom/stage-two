package utils

import (
	"os"
	"time"

	"github.com/Micah-Shallom/stage-two/models"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type SignedDetails struct {
	FirstName string
	LastName  string
	Email     string
	UserID    string
	UserType  string
	jwt.RegisteredClaims
}

func GenerateJWT(user *models.User) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour) // 1 week
	claims := &SignedDetails{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		UserID:    user.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.UserID,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return ss, nil
}
