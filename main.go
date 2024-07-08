package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Micah-Shallom/stage-two/config"
	"github.com/Micah-Shallom/stage-two/handlers"
	"github.com/Micah-Shallom/stage-two/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	//load dotenv
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Error loading .env file")
	}

	//initialize config
	conf := config.NewConfig()

	//initialize logger
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	//initialize database
	db, err := config.OpenDB()
	if err != nil {
		logger.Fatal(err)
	}

	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Failed to get postgresDB from GORM: %v", err)
		}
		sqlDB.Close()
	}()

	_ = conf.Database.Port
	dsn := conf.Database.DSN

	//start the server
	logger.Printf("Starting server on port %s", dsn)

	//setup gin router
	//setup a new gin router
	r := gin.Default()

	//initialize application
	app := config.NewApplication(db)
	//initialize handlers
	handler := handlers.NewHandlers(app)
	//initialize routes and pass the router and handler
	router := routes.Routes(r, handler)



	portStr := fmt.Sprintf(":8080")
	fmt.Println(portStr, dsn)
	router.Run(portStr)
}
