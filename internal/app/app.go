package app

import (
	"fmt"
	"github.com/helmigandi/go-workout-api/internal/api"
	"log"
	"net/http"
	"os"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
}

// NewApplication creates a formatted print line across the application.
func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// our stores will go here

	// our handlers will go here
	workoutHandler := api.NewWorkoutHandler()

	return &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
	}, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status available\n")
}
