package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/maliByatzes/fwt"
	"github.com/maliByatzes/fwt/postgres"
	"github.com/stretchr/testify/require"
)

func TestWorkoutService_CreateWorkout(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)
		s := postgres.NewWorkoutService(db)

		user, ctx := MustCreateUser(t, context.Background(), db, &fwt.User{
			Username:       postgres.RandomUsername(),
			Email:          postgres.RandomEmail(),
			HashedPassword: postgres.RandomHashedPassword(),
		})

		exercise1 := MustCreateExercise(t, ctx, db, &fwt.Exercise{Name: postgres.RandomString(12), Description: postgres.RandomString(50)})
		exercise2 := MustCreateExercise(t, ctx, db, &fwt.Exercise{Name: postgres.RandomString(12), Description: postgres.RandomString(50)})
		exercise3 := MustCreateExercise(t, ctx, db, &fwt.Exercise{Name: postgres.RandomString(12), Description: postgres.RandomString(50)})

		newWorkout := &fwt.Workout{
			UserID:        user.ID,
			Name:          postgres.RandomString(12),
			ScheduledDate: time.Now().Add(time.Hour),
			Exercises:     []*fwt.Exercise{exercise1, exercise2, exercise3},
		}

		err := s.CreateWorkout(ctx, newWorkout)
		require.NoError(t, err)

		require.Equal(t, newWorkout.ID, uint(1))
		require.NotZero(t, newWorkout.CreatedAt)
		require.NotZero(t, newWorkout.UpdatedAt)
	})

	t.Run("ErrNameRequired", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)
		s := postgres.NewWorkoutService(db)

		user, ctx := MustCreateUser(t, context.Background(), db, &fwt.User{
			Username:       postgres.RandomUsername(),
			Email:          postgres.RandomEmail(),
			HashedPassword: postgres.RandomHashedPassword(),
		})

		exercise1 := MustCreateExercise(t, ctx, db, &fwt.Exercise{Name: postgres.RandomString(12), Description: postgres.RandomString(50)})
		exercise2 := MustCreateExercise(t, ctx, db, &fwt.Exercise{Name: postgres.RandomString(12), Description: postgres.RandomString(50)})
		exercise3 := MustCreateExercise(t, ctx, db, &fwt.Exercise{Name: postgres.RandomString(12), Description: postgres.RandomString(50)})

		newWorkout := &fwt.Workout{
			UserID:        user.ID,
			ScheduledDate: time.Now().Add(time.Hour),
			Exercises:     []*fwt.Exercise{exercise1, exercise2, exercise3},
		}

		err := s.CreateWorkout(ctx, newWorkout)
		require.Error(t, err)
		require.Equal(t, fwt.ErrorCode(err), fwt.EINVALID)
		require.Equal(t, fwt.ErrorMessage(err), "Name is required.")
	})

	t.Run("ErrSDRequired", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)
		s := postgres.NewWorkoutService(db)

		user, ctx := MustCreateUser(t, context.Background(), db, &fwt.User{
			Username:       postgres.RandomUsername(),
			Email:          postgres.RandomEmail(),
			HashedPassword: postgres.RandomHashedPassword(),
		})

		exercise1 := MustCreateExercise(t, ctx, db, &fwt.Exercise{Name: postgres.RandomString(12), Description: postgres.RandomString(50)})
		exercise2 := MustCreateExercise(t, ctx, db, &fwt.Exercise{Name: postgres.RandomString(12), Description: postgres.RandomString(50)})
		exercise3 := MustCreateExercise(t, ctx, db, &fwt.Exercise{Name: postgres.RandomString(12), Description: postgres.RandomString(50)})

		newWorkout := &fwt.Workout{
			UserID:    user.ID,
			Name:      postgres.RandomString(12),
			Exercises: []*fwt.Exercise{exercise1, exercise2, exercise3},
		}

		err := s.CreateWorkout(ctx, newWorkout)
		require.Error(t, err)
		require.Equal(t, fwt.ErrorCode(err), fwt.EINVALID)
		require.Equal(t, fwt.ErrorMessage(err), "Scheduled Date is required.")
	})
}

func TestWorkoutService_FindWorkouts(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)
		s := postgres.NewWorkoutService(db)

		user, ctx := MustCreateUser(t, context.Background(), db, &fwt.User{Username: postgres.RandomUsername(), Email: postgres.RandomEmail(), HashedPassword: postgres.RandomHashedPassword()})

		exercise1 := MustCreateExercise(t, ctx, db, &fwt.Exercise{Name: postgres.RandomString(12), Description: postgres.RandomString(50)})
		exercise2 := MustCreateExercise(t, ctx, db, &fwt.Exercise{Name: postgres.RandomString(12), Description: postgres.RandomString(50)})
		exercise3 := MustCreateExercise(t, ctx, db, &fwt.Exercise{Name: postgres.RandomString(12), Description: postgres.RandomString(50)})
		MustCreateWorkout(t, ctx, db, &fwt.Workout{
			UserID:        user.ID,
			Name:          postgres.RandomString(12),
			ScheduledDate: time.Now().Add(time.Hour),
			Exercises:     []*fwt.Exercise{exercise1, exercise2, exercise3},
		})

		id := uint(1)
		a, n, err := s.FindWorkouts(ctx, fwt.WorkoutFilter{ID: &id})
		require.NoError(t, err)
		require.Equal(t, len(a), 1)
		require.Equal(t, a[0].ID, id)
		require.Equal(t, n, 3)
	})
}

func MustCreateWorkout(tb testing.TB, ctx context.Context, db *postgres.DB, workout *fwt.Workout) *fwt.Workout {
	tb.Helper()
	err := postgres.NewWorkoutService(db).CreateWorkout(ctx, workout)
	require.NoError(tb, err)
	return workout
}
