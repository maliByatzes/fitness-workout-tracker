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

		ctx := context.Background()
		user := MustCreateUser(t, ctx, db, &fwt.User{
			Username:       postgres.RandomUsername(),
			Email:          postgres.RandomEmail(),
			HashedPassword: postgres.RandomHashedPassword(),
		})

		newWorkout := &fwt.Workout{
			UserID:        user.ID,
			Name:          postgres.RandomString(12),
			ScheduledDate: time.Now().Add(time.Hour),
		}

		err := s.CreateWorkout(ctx, newWorkout)
		require.NoError(t, err)

		require.Equal(t, newWorkout.ID, uint(1))
		require.NotZero(t, newWorkout.CreatedAt)
		require.NotZero(t, newWorkout.UpdatedAt)
	})
}
