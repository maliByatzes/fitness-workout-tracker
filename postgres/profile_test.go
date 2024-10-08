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

		user, ctx := MustCreateUser(t, context.Background(), db, &fwt.User{
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
		require.Equal(t, fwt.ErrorCode(err), fwt.ENOTAUTHORIZED)
		require.Equal(t, fwt.ErrorMessage(err), "You must be logged in to create a profile.")
	})
}

func TestProfileService_FindProfiles(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)
		s := postgres.NewProfileService(db)

		user, ctx := MustCreateUser(t, context.Background(), db, &fwt.User{Username: "kyle", Email: "kyle@email.com", HashedPassword: "password"})

		MustCreateProfile(t, ctx, db, &fwt.Profile{UserID: user.ID, FirstName: "kyle1"})

		id := uint(1)
		a, n, err := s.FindProfiles(ctx, fwt.ProfileFilter{ID: &id})
		require.NoError(t, err)
		require.Equal(t, len(a), 1)
		require.Equal(t, a[0].ID, id)
		require.Equal(t, a[0].FirstName, "kyle1")
		require.Equal(t, n, 1)
	})
}

func TestProfileService_UpdateProfile(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)
		s := postgres.NewProfileService(db)
		user, ctx := MustCreateUser(t, context.Background(), db, &fwt.User{
			Username:       "jeff",
			Email:          "jeff@email.com",
			HashedPassword: "password",
		})
		profile0 := MustCreateProfile(t, ctx, db, &fwt.Profile{
			UserID:    user.ID,
			FirstName: "jeffina",
		})

		newFirstName := "kyle"
		up, err := s.UpdateProfile(ctx, profile0.ID, fwt.ProfileUpdate{
			FirstName: &newFirstName,
		})
		require.NoError(t, err)
		require.Equal(t, up.FirstName, newFirstName)

		other, err := s.FindProfileByID(ctx, 1)
		require.NoError(t, err)
		require.Equal(t, up, other)
	})
}

func TestProfileService_DeleteProfile(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)
	s := postgres.NewProfileService(db)
	user0, ctx0 := MustCreateUser(t, context.Background(), db, &fwt.User{
		Username:       "jeff",
		Email:          "jeff@email.com",
		HashedPassword: "password",
	})
	profile0 := MustCreateProfile(t, ctx0, db, &fwt.Profile{
		UserID:    user0.ID,
		FirstName: "jeffina",
		LastName:  "reboot",
	})

	err := s.DeleteProfile(ctx0, profile0.ID)
	require.NoError(t, err)

	_, err = s.FindProfileByID(ctx0, profile0.ID)
	require.Error(t, err)
	require.Equal(t, fwt.ErrorCode(err), fwt.ENOTFOUND)
	require.Equal(t, fwt.ErrorMessage(err), "Profile not found.")
}

func MustCreateProfile(tb testing.TB, ctx context.Context, db *postgres.DB, profile *fwt.Profile) *fwt.Profile {
	tb.Helper()
	err := postgres.NewProfileService(db).CreateProfile(ctx, profile)
	require.NoError(tb, err)
	return profile
}
