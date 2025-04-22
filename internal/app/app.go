package app

import (
	"log"
	"os"
)

type Application struct {
	Logger *log.Logger
}

// NewApplication creates a formatted print line across the application.
func NewApplication() (*Application, error) {
	return &Application{
		Logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}, nil
}
