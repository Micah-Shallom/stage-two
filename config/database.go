package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Micah-Shallom/stage-two/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Port         string
	Env          string
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

func NewDatabase() *DatabaseConfig {
	return &DatabaseConfig{
		Env:          os.Getenv("ENV"),
		DSN:          os.Getenv("DB_DSN"),
		Port:         os.Getenv("PORT"),
		MaxOpenConns: 25,
		MaxIdleConns: 25,
		MaxIdleTime:  os.Getenv("DB_MAX_IDLE_TIME"),
	}
}

func OpenDB() (*gorm.DB, error) {
	cfg := NewDatabase()

	dsn := cfg.DSN
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	//run auto migration
	err = db.AutoMigrate(
		&models.User{},
		&models.Organisation{},
	)

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

	duration, err := time.ParseDuration(cfg.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	sqlDB.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = sqlDB.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	logger := log.New(log.Writer(), "", log.Ldate|log.Ltime)

	logger.Println("Database connection successful")

	return db, nil
}
