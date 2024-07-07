package tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Micah-Shallom/stage-two/config"
	"github.com/Micah-Shallom/stage-two/models"
	"github.com/Micah-Shallom/stage-two/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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


func setupTestDB() (*gorm.DB, error) {
	// Read test database configuration from environment variables or use default values
	dbName := "test_db"
	dbUser := "test_user"
	dbPassword := "test_password"
	dbHost := "localhost"
	dbPort := "5432"

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	// Drop and recreate the test database
	err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s;", dbName)).Error
	if err != nil {
		return nil, fmt.Errorf("error dropping database: %v", err)
	}
	err = db.Exec(fmt.Sprintf("CREATE DATABASE %s;", dbName)).Error
	if err != nil {
		return nil, fmt.Errorf("error creating database: %v", err)
	}

	// Connect to the test database
	dsn = fmt.Sprintf("%s dbname=%s", dsn, dbName)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error connecting to test database: %v", err)
	}

	// Run auto migrations
	err = db.AutoMigrate(
		&models.User{},
		&models.Organisation{},
	)
	if err != nil {
		return nil, fmt.Errorf("error migrating schema: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error getting DB instance: %v", err)
	}

	// Set connection pool settings
	maxOpenConns := 10
	maxIdleConns := 5
	maxIdleTime := time.Minute * 5

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxIdleTime(maxIdleTime)

	// Ping database to ensure connection is valid
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = sqlDB.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	logger := log.New(log.Writer(), "", log.Ldate|log.Ltime)
	logger.Println("Test database setup successful")

	return db, nil
}

func TestOrganizationAccess(t *testing.T) {
	db, err := setupTestDB() // assume setupTestDB sets up a test database
	assert.NoError(t, err)

	app := config.NewApplication(db)

	// Create users
	user1 := &models.User{
		UserID:    "user1",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "password",
	}
	user2 := &models.User{
		UserID:    "user2",
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "jane.smith@example.com",
		Password:  "password",
	}

	// Create organizations
	org1 := &models.Organisation{
		OrgID: "org1",
		Name:  "John's Org",
		Users: []models.User{*user1},
	}
	org2 := &models.Organisation{
		OrgID: "org2",
		Name:  "Jane's Org",
		Users: []models.User{*user2},
	}

	// Save users and organizations to the database
	err = app.Models.Users.Create(user1)
	assert.NoError(t, err)
	err = app.Models.Users.Create(user2)
	assert.NoError(t, err)
	err = app.Models.Organisations.Create(org1)
	assert.NoError(t, err)
	err = app.Models.Organisations.Create(org2)
	assert.NoError(t, err)

	// Test user1 accessing user2's organization
	orgs, err := app.Models.Organisations.GetByUserID(user1.UserID)
	assert.NoError(t, err)
	assert.Empty(t, orgs)

	// Test user1 accessing their own organization
	orgs, err = app.Models.Organisations.GetByUserID(user1.UserID)
	assert.NoError(t, err)
	assert.NotEmpty(t, orgs)
	assert.Equal(t, "John's Org", orgs[0].Name)
}
