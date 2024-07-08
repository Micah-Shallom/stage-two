package tests

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/Micah-Shallom/stage-two/models"
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

// ConnectDB - It will be used anywhere in the application to connect database.
func SetupTestDB() (*gorm.DB, error) {
	test_dsn := os.Getenv("TEST_DSN")

	db, err := gorm.Open(postgres.Open(test_dsn), &gorm.Config{})
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

	return db, nil
}

// MockedDB is used in unit tests to mock db
func MockedDB(operation string) {
	/*
	   If tests are running in CI, environment variables should not be loaded.
	   The reason is environment vars will be provided through CI config file.
	*/
	if CI := os.Getenv("CI"); CI == "" {
		// If tests are not running in CI, we have to load .env file.
		_, fileName, _, _ := runtime.Caller(0)
		currPath := filepath.Dir(fileName)
		// path should be relative path from this directory to ".env"
		err := godotenv.Load(currPath + "/../../.env")
		if err != nil {
			log.Fatalf("Error loading env.\n%v", err)
		}
	}

	dbName := os.Getenv("DATABASE_NAME")
	pgUser := os.Getenv("PSQL_USER")
	pgPassword := os.Getenv("PSQL_PASSWORD")

	// createdb => https://www.postgresql.org/docs/7.0/app-createdb.htm
	// dropdb => https://www.postgresql.org/docs/7.0/app-dropdb.htm
	var command string

	if operation == CREATE {
		command = "createdb"
	} else {
		command = "dropdb"
	}

	// createdb & dropdb commands have same configuration syntax.
	cmd := exec.Command(command, "-h", "localhost", "-U", pgUser, "-e", dbName)
	cmd.Env = os.Environ()

	/*
	   if we normally execute createdb/dropdb, we will be prompted to provide password.
	   To inject password automatically, we have to set PGPASSWORD as prefix.
	*/
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%v", pgPassword))

	if err := cmd.Run(); err != nil {
		log.Fatalf("Error executing %v on %v.\n%v", command, dbName, err)
	}

	/*
	   Alternatively instead of createdb/dropdb, you can use
	   psql -c "CREATE/DROP DATABASE DBNAME" "DATABASE_URL"
	*/
}
