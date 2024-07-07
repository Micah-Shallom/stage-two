package config

import (
	"log"
	"os"

	"github.com/Micah-Shallom/stage-two/models"
	"gorm.io/gorm"
)

type Application struct {
	Config *Config
	Logger *log.Logger
	Models models.Models
}

func NewApplication(db *gorm.DB) *Application {
	cfg := NewConfig()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	return &Application{
		Config: cfg,
		Logger: logger,
		Models: models.NewModels(db),
	}
}
