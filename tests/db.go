package tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Micah-Shallom/stage-two/models"
	"github.com/go-pg/pg/v10"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	// CREATE - operation to create database
	CREATE = "CREATE"
	// DROP - operation to drop database
	DROP = "DROP"
)

func connectDB() (*gorm.DB, error) {
	test_dsn := os.Getenv("TEST_DSN")

	db, err := gorm.Open(postgres.Open(test_dsn), &gorm.Config{})
	fmt.Println(err)

	if err != nil {
		log.Printf("Error connecting database.\n%v", err)
		return nil, err
	}

	// Run auto migration
	err = db.AutoMigrate(
		&models.User{},
		&models.Organisation{},
	)
	if err != nil {
		return nil, err
	}

	log.Println("Database Connected Successfully")

	return db, nil
}

// ConnectDB - It will be used anywhere in the application to connect database.
func SetupTestDB() (*gorm.DB, error) {
	// Create the database
	MockedDB(CREATE)

	//Retry mechanism
	for i := 0; i < 5; i++ {
		db, err := connectDB()
		if err != nil {
			log.Printf("Failed to connect to the database: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		return db, nil
	}

	return nil, fmt.Errorf("Failed to connect to the database after 5 retries")
}

// MockedDB is used in unit tests to mock db
func MockedDB(operation string) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading env %v", err)
	}

	pgUser := os.Getenv("PSQL_USER")
	pgPassword := os.Getenv("PSQL_PASSWORD")
	dbName := "testdb"

	options := &pg.Options{
		User:     pgUser,
		Password: pgPassword,
		Addr:     "localhost:5432",
		Database: "postgres", // Connect to the postgres database first
	}

	db := pg.Connect(options)
	defer db.Close()

	ctx := context.Background()

	if operation == CREATE {

		_, err := db.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			fmt.Printf("Error creating database %s: %v", dbName, err)
			return
		}
		log.Printf("Database %s created successfully", dbName)
	} else if operation == DROP {
		// Terminate all connections to the database
		_, err := db.ExecContext(ctx, fmt.Sprintf(`
			SELECT pg_terminate_backend(pg_stat_activity.pid)
			FROM pg_stat_activity
			WHERE pg_stat_activity.datname = '%s' AND pid <> pg_backend_pid();
		`, dbName))
		if err != nil {
			fmt.Printf("Error terminating connections to database %s: %v", dbName, err)
			return
		}

		_, err = db.ExecContext(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
		if err != nil {
			fmt.Printf("Error dropping database %s: %v", dbName, err)
			return
		}
		log.Printf("Database %s dropped successfully", dbName)
	}
}
