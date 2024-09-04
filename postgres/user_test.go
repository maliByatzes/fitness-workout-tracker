package postgres_test

import (
	"context"
	"testing"

	"github.com/maliByatzes/fwt"
	"github.com/maliByatzes/fwt/postgres"
	"github.com/stretchr/testify/require"
)

func TestUserService_CreateUser(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseBD(t, db)
		s := postgres.NewUserService(db)

		newUser := &fwt.User{
			Username:       "jane",
			Email:          "jane@email.com",
			HashedPassword: "hashed_password",
		}

		err := s.CreateUser(context.Background(), newUser)
		require.NoError(t, err)

		got, want := newUser.ID, uint(1)
		require.Equal(t, got, want)
		require.False(t, newUser.CreatedAt.IsZero())
		require.False(t, newUser.UpdatedAt.IsZero())
	})

	t.Run("ErrUsernameRequired", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseBD(t, db)
		s := postgres.NewUserService(db)
		err := s.CreateUser(context.Background(), &fwt.User{})
		require.Error(t, err)
		require.Equal(t, fwt.ErrorCode(err), fwt.EINVALID)
		require.Equal(t, fwt.ErrorMessage(err), "Username is required.")
	})

	t.Run("ErrEmailRequired", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseBD(t, db)
		s := postgres.NewUserService(db)
		err := s.CreateUser(context.Background(), &fwt.User{Username: "jane"})
		require.Error(t, err)
		require.Equal(t, fwt.ErrorCode(err), fwt.EINVALID)
		require.Equal(t, fwt.ErrorMessage(err), "Email is required.")
	})
}
