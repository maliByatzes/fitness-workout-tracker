package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/maliByatzes/fwt"
	"github.com/maliByatzes/fwt/postgres"
	"github.com/stretchr/testify/require"
)

func TestWorkoutExerciseService_CreateWorkoutExercise(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)
		s := postgres.NewWorkoutExerciseService(db)

		user, ctx := MustCreateUser(t, context.Background(), db, &fwt.User{Username: postgres.RandomUsername(), Email: postgres.RandomEmail(), HashedPassword: postgres.RandomHashedPassword()})
		workout := MustCreateWorkout(t, ctx, db, &fwt.Workout{UserID: user.ID, Name: postgres.RandomString(6), ScheduledDate: time.Now().Add(time.Hour)})
		exercise := MustCreateExercise(t, ctx, db, &fwt.Exercise{Name: postgres.RandomString(12), Description: postgres.RandomString(50)})

		newWorkoutExercise := &fwt.WorkoutExercise{
			WorkoutID:  workout.ID,
			ExerciseID: exercise.ID,
			Order:      1,
		}

		err := s.CreateWorkoutExercise(ctx, newWorkoutExercise)
		require.NoError(t, err)
		require.Equal(t, newWorkoutExercise.ID, uint(4))
		require.NotZero(t, newWorkoutExercise.CreatedAt)
		require.NotZero(t, newWorkoutExercise.UpdatedAt)
	})

	t.Run("ErrWorkoutIDRequired", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)
		s := postgres.NewWorkoutExerciseService(db)

		ctx := context.Background()
		exercise := MustCreateExercise(t, ctx, db, &fwt.Exercise{Name: postgres.RandomString(12), Description: postgres.RandomString(50)})

		newWorkoutExercise := &fwt.WorkoutExercise{
			ExerciseID: exercise.ID,
			Order:      1,
		}

		err := s.CreateWorkoutExercise(ctx, newWorkoutExercise)
		require.Error(t, err)
		require.Equal(t, fwt.ErrorCode(err), fwt.EINVALID)
		require.Equal(t, fwt.ErrorMessage(err), "WorkoutID is required.")
	})

	t.Run("ErrExerciseIDRequired", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)
		s := postgres.NewWorkoutExerciseService(db)

		user, ctx := MustCreateUser(t, context.Background(), db, &fwt.User{Username: postgres.RandomUsername(), Email: postgres.RandomEmail(), HashedPassword: postgres.RandomHashedPassword()})
		workout := MustCreateWorkout(t, ctx, db, &fwt.Workout{UserID: user.ID, Name: postgres.RandomString(6), ScheduledDate: time.Now().Add(time.Hour)})

		newWorkoutExercise := &fwt.WorkoutExercise{
			WorkoutID: workout.ID,
			Order:     1,
		}

		err := s.CreateWorkoutExercise(ctx, newWorkoutExercise)
		require.Error(t, err)
		require.Equal(t, fwt.ErrorCode(err), fwt.EINVALID)
		require.Equal(t, fwt.ErrorMessage(err), "ExerciseID is required.")
	})

	t.Run("ErrOrderRequired", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)
		s := postgres.NewWorkoutExerciseService(db)

		user, ctx := MustCreateUser(t, context.Background(), db, &fwt.User{Username: postgres.RandomUsername(), Email: postgres.RandomEmail(), HashedPassword: postgres.RandomHashedPassword()})
		workout := MustCreateWorkout(t, ctx, db, &fwt.Workout{UserID: user.ID, Name: postgres.RandomString(6), ScheduledDate: time.Now().Add(time.Hour)})
		exercise := MustCreateExercise(t, ctx, db, &fwt.Exercise{Name: postgres.RandomString(12), Description: postgres.RandomString(50)})

		newWorkoutExercise := &fwt.WorkoutExercise{
			WorkoutID:  workout.ID,
			ExerciseID: exercise.ID,
		}

		err := s.CreateWorkoutExercise(ctx, newWorkoutExercise)
		require.Error(t, err)
		require.Equal(t, fwt.ErrorCode(err), fwt.EINVALID)
		require.Equal(t, fwt.ErrorMessage(err), "Order is required.")
	})
}
