package store

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
	if err != nil {
		t.Fatalf("opening test db: %v", err)
	}

	// run the migratoins for our test db
	err = Migrate(db, "../../migrations/")
	if err != nil {
		t.Fatalf("migrating test db error: %v", err)
	}

	// wipe the DB every time if we want to run our setup tests.
	_, err = db.Exec(`TRUNCATE users, workouts, workout_entries CASCADE`)
	if err != nil {
		t.Fatalf("truncating tables %v", err)
	}

	return db
}

func TestCreateWorkout(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewPostgresWorkoutStore(db)
	tests := []struct {
		name    string
		workout *Workout
		wantErr bool
	}{
		{
			name: "valid workout",
			workout: &Workout{
				Title:           "push day",
				Description:     "upper body day",
				DurationMinutes: 60,
				CaloriesBurned:  200,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Bench Press",
						Sets:         3,
						Reps:         createIntPtr(10),
						Weight:       createFloatPtr(85.9),
						Notes:        "warm up properly",
						OrderIndex:   1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "workout with invalid entries",
			workout: &Workout{
				Title:           "full body",
				Description:     "full body day",
				DurationMinutes: 90,
				CaloriesBurned:  500,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Plank",
						Sets:         3,
						Reps:         createIntPtr(60),
						Notes:        "keep form",
						OrderIndex:   1,
					},
					{
						ExerciseName:    "Squats",
						Sets:            4,
						Reps:            createIntPtr(12),
						DurationSeconds: createIntPtr(60),
						Weight:          createFloatPtr(101.1),
						OrderIndex:      2,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := store.CreateWorkout(tt.workout)
			// without assertions
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateWorkout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// with assertions
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.workout.Title, got.Title)
			assert.Equal(t, tt.workout.Description, got.Description)
			assert.Equal(t, tt.workout.DurationMinutes, got.DurationMinutes)
			assert.Equal(t, tt.workout.CaloriesBurned, got.CaloriesBurned)
			assert.Equal(t, len(tt.workout.Entries), len(got.Entries))

			savedWorkout, err := store.GetWorkoutByID(int64(got.ID))
			require.NoError(t, err)
			assert.Equal(t, got.ID, savedWorkout.ID)
			assert.Equal(t, got.Title, savedWorkout.Title)
			assert.Equal(t, got.Description, savedWorkout.Description)
			assert.Equal(t, got.DurationMinutes, savedWorkout.DurationMinutes)
			assert.Equal(t, got.CaloriesBurned, savedWorkout.CaloriesBurned)
			assert.Equal(t, len(got.Entries), len(savedWorkout.Entries))

			for i, entry := range got.Entries {
				assert.Equal(t, entry, savedWorkout.Entries[i])
				assert.Equal(t, entry.OrderIndex, savedWorkout.Entries[i].OrderIndex)
				assert.Equal(t, entry.ExerciseName, savedWorkout.Entries[i].ExerciseName)
				assert.Equal(t, entry.Sets, savedWorkout.Entries[i].Sets)
				assert.Equal(t, entry.Reps, savedWorkout.Entries[i].Reps)
				assert.Equal(t, entry.Weight, savedWorkout.Entries[i].Weight)
				assert.Equal(t, entry.Notes, savedWorkout.Entries[i].Notes)
				assert.Equal(t, entry.DurationSeconds, savedWorkout.Entries[i].DurationSeconds)
			}
		})
	}
}

func createIntPtr(i int) *int {
	return &i
}

func createFloatPtr(f float64) *float64 {
	return &f
}
