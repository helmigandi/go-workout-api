package app

import (
	"database/sql"
	"fmt"
	"github.com/helmigandi/go-workout-api/internal/api"
	"github.com/helmigandi/go-workout-api/internal/store"
	"github.com/helmigandi/go-workout-api/migrations"
	"log"
	"net/http"
	"os"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	DB             *sql.DB
}

// NewApplication creates a formatted print line across the application.
func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}
	
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// our stores will go here

	// our handlers will go here
	workoutHandler := api.NewWorkoutHandler()

	return &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		DB:             pgDB,
	}, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status available\n")
}
