package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/maliByatzes/fwt"
	"github.com/maliByatzes/fwt/postgres"
	"github.com/stretchr/testify/require"
)

func TestProfileService_CreateProfile(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)
		s := postgres.NewProfileService(db)

		ctx := context.Background()
		user := MustCreateUser(t, ctx, db, &fwt.User{
			Username:       "jeff",
			Email:          "jeff@email.com",
			HashedPassword: "password",
		})

		newProfile := &fwt.Profile{
			UserID:      user.ID,
			FirstName:   "jeffina",
			LastName:    "robertson",
			DateOfBirth: time.Now(),
			Gender:      "Male",
			Height:      float64(167.34),
			Weight:      float64(55.4),
		}

		err := s.CreateProfile(ctx, newProfile)
		require.NoError(t, err)

		got, want := newProfile.ID, uint(1)
		require.Equal(t, got, want)
		require.NotZero(t, newProfile.CreatedAt)
		require.NotZero(t, newProfile.UpdatedAt)
	})

	t.Run("ErrUserIDRequired", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)
		s := postgres.NewProfileService(db)
		err := s.CreateProfile(context.Background(), &fwt.Profile{})
		require.Error(t, err)
		require.Equal(t, fwt.ErrorCode(err), fwt.EINVALID)
		require.Equal(t, fwt.ErrorMessage(err), "UserID is required.")
	})
}
