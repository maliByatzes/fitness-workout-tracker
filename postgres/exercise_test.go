package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/maliByatzes/fwt"
	"github.com/maliByatzes/fwt/postgres"
	"github.com/stretchr/testify/require"
)

func TestExerciseService_CreateExercise(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)
		s := postgres.NewExerciseService(db)

		newExercise := &fwt.Exercise{
			Name:        postgres.RandomString(12),
			Description: postgres.RandomString(50),
		}

		err := s.CreateExercise(context.Background(), newExercise)
		require.NoError(t, err)

		fmt.Println("newExercise.ID in OK: ", newExercise.ID)
		require.Equal(t, newExercise.ID, uint(1))
		require.NotZero(t, newExercise.CreatedAt)
		require.NotZero(t, newExercise.UpdatedAt)
	})

	t.Run("ErrNameRequired", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)
		s := postgres.NewExerciseService(db)

		newExercise := &fwt.Exercise{
			Description: postgres.RandomString(50),
		}
		err := s.CreateExercise(context.Background(), newExercise)
		require.Error(t, err)
		require.Equal(t, fwt.ErrorCode(err), fwt.EINVALID)
		require.Equal(t, fwt.ErrorMessage(err), "Name is required.")
	})

	t.Run("ErrDescriptionRequired", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)
		s := postgres.NewExerciseService(db)

		newExercise := &fwt.Exercise{
			Name: postgres.RandomString(12),
		}
		err := s.CreateExercise(context.Background(), newExercise)
		require.Error(t, err)
		require.Equal(t, fwt.ErrorCode(err), fwt.EINVALID)
		require.Equal(t, fwt.ErrorMessage(err), "Description is required.")
	})
}

func TestExerciseService_FindExercises(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)
	s := postgres.NewExerciseService(db)

	ctx := context.Background()
	MustCreateExercise(t, ctx, db, &fwt.Exercise{
		Name:        postgres.RandomString(12),
		Description: postgres.RandomString(50),
	})

	id := uint(1)
	a, n, err := s.FindExercises(ctx, fwt.ExerciseFilter{ID: &id})
	require.NoError(t, err)
	require.Equal(t, len(a), 1)
	require.Equal(t, a[0].ID, id)
	require.Equal(t, n, 1)
}

func MustCreateExercise(tb testing.TB, ctx context.Context, db *postgres.DB, exercise *fwt.Exercise) *fwt.Exercise {
	tb.Helper()
	err := postgres.NewExerciseService(db).CreateExercise(ctx, exercise)
	require.NoError(tb, err)
	return exercise
}
