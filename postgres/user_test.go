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
		defer MustCloseDB(t, db)
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
		defer MustCloseDB(t, db)
		s := postgres.NewUserService(db)
		err := s.CreateUser(context.Background(), &fwt.User{})
		require.Error(t, err)
		require.Equal(t, fwt.ErrorCode(err), fwt.EINVALID)
		require.Equal(t, fwt.ErrorMessage(err), "Username is required.")
	})

	t.Run("ErrEmailRequired", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)
		s := postgres.NewUserService(db)
		err := s.CreateUser(context.Background(), &fwt.User{Username: "jane"})
		require.Error(t, err)
		require.Equal(t, fwt.ErrorCode(err), fwt.EINVALID)
		require.Equal(t, fwt.ErrorMessage(err), "Email is required.")
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)
		s := postgres.NewUserService(db)
		user0 := MustCreateUser(t, context.Background(), db, &fwt.User{
			Username:       "janedoe",
			Email:          "janedoe@email.com",
			HashedPassword: "password",
		})

		newUsername, newEmail := "jill", "jill@email.com"
		uu, err := s.UpdateUser(context.Background(), user0.ID, fwt.UserUpdate{
			Username: &newUsername,
			Email:    &newEmail,
		})
		require.NoError(t, err)
		require.Equal(t, uu.Username, newUsername)
		require.Equal(t, uu.Email, newEmail)

		other, err := s.FindUserByID(context.Background(), 1)
		require.NoError(t, err)
		require.Equal(t, uu, other)
	})

	t.Run("UpdateUsername", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)
		s := postgres.NewUserService(db)
		user0 := MustCreateUser(t, context.Background(), db, &fwt.User{
			Username:       "janedoe",
			Email:          "janedoe@email.com",
			HashedPassword: "password",
		})

		newUsername := "jill"
		uu, err := s.UpdateUser(context.Background(), user0.ID, fwt.UserUpdate{
			Username: &newUsername,
		})
		require.NoError(t, err)
		require.Equal(t, uu.Username, newUsername)

		other, err := s.FindUserByID(context.Background(), 1)
		require.NoError(t, err)
		require.Equal(t, uu, other)
	})

	t.Run("UpdateEmail", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)
		s := postgres.NewUserService(db)
		user0 := MustCreateUser(t, context.Background(), db, &fwt.User{
			Username:       "janedoe",
			Email:          "janedoe@email.com",
			HashedPassword: "password",
		})

		newEmail := "jill@email.com"
		uu, err := s.UpdateUser(context.Background(), user0.ID, fwt.UserUpdate{
			Email: &newEmail,
		})
		require.NoError(t, err)
		require.Equal(t, uu.Email, newEmail)

		other, err := s.FindUserByID(context.Background(), 1)
		require.NoError(t, err)
		require.Equal(t, uu, other)
	})

	t.Run("UpdateNothing", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)
		s := postgres.NewUserService(db)
		user0 := MustCreateUser(t, context.Background(), db, &fwt.User{
			Username:       "janedoe",
			Email:          "janedoe@email.com",
			HashedPassword: "password",
		})

		uu, err := s.UpdateUser(context.Background(), user0.ID, fwt.UserUpdate{})
		require.NoError(t, err)

		other, err := s.FindUserByID(context.Background(), 1)
		require.NoError(t, err)
		require.Equal(t, uu, other)
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)
		s := postgres.NewUserService(db)
		user0 := MustCreateUser(t, context.Background(), db, &fwt.User{
			Username:       "jeff",
			Email:          "jeff@email.com",
			HashedPassword: "password",
		})

		err := s.DeleteUser(context.Background(), user0.ID)
		require.NoError(t, err)
	})
}

func TestUserService_FindUsers(t *testing.T) {
	t.Run("ID", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)
		s := postgres.NewUserService(db)

		ctx := context.Background()
		MustCreateUser(t, ctx, db, &fwt.User{Username: "janedoe", Email: "janedoe@email.com", HashedPassword: "password"})
		MustCreateUser(t, ctx, db, &fwt.User{Username: "kyledoe", Email: "kyledoe@email.com", HashedPassword: "password"})
		MustCreateUser(t, ctx, db, &fwt.User{Username: "jimdoe", Email: "jimdoe@email.com", HashedPassword: "password"})
		MustCreateUser(t, ctx, db, &fwt.User{Username: "frankdoe", Email: "frankdoe@email.com", HashedPassword: "password"})

		id := uint(2)
		a, n, err := s.FindUsers(ctx, fwt.UserFilter{ID: &id})
		require.NoError(t, err)
		require.Equal(t, len(a), 1)
		require.Equal(t, a[0].ID, uint(2))
		require.Equal(t, a[0].Username, "kyledoe")
		require.Equal(t, a[0].Email, "kyledoe@email.com")
		require.Equal(t, n, 1)
	})

	t.Run("Username", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)
		s := postgres.NewUserService(db)

		ctx := context.Background()
		MustCreateUser(t, ctx, db, &fwt.User{Username: "janedoe", Email: "janedoe@email.com", HashedPassword: "password"})
		MustCreateUser(t, ctx, db, &fwt.User{Username: "kyledoe", Email: "kyledoe@email.com", HashedPassword: "password"})
		MustCreateUser(t, ctx, db, &fwt.User{Username: "jimdoe", Email: "jimdoe@email.com", HashedPassword: "password"})
		MustCreateUser(t, ctx, db, &fwt.User{Username: "frankdoe", Email: "frankdoe@email.com", HashedPassword: "password"})

		username := "frankdoe"
		a, n, err := s.FindUsers(ctx, fwt.UserFilter{Username: &username})
		require.NoError(t, err)
		require.Equal(t, len(a), 1)
		require.Equal(t, a[0].ID, uint(4))
		require.Equal(t, a[0].Username, username)
		require.Equal(t, a[0].Email, "frankdoe@email.com")
		require.Equal(t, n, 1)
	})

	t.Run("Email", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)
		s := postgres.NewUserService(db)

		ctx := context.Background()
		MustCreateUser(t, ctx, db, &fwt.User{Username: "janedoe", Email: "janedoe@email.com", HashedPassword: "password"})
		MustCreateUser(t, ctx, db, &fwt.User{Username: "kyledoe", Email: "kyledoe@email.com", HashedPassword: "password"})
		MustCreateUser(t, ctx, db, &fwt.User{Username: "jimdoe", Email: "jimdoe@email.com", HashedPassword: "password"})
		MustCreateUser(t, ctx, db, &fwt.User{Username: "frankdoe", Email: "frankdoe@email.com", HashedPassword: "password"})

		email := "janedoe@email.com"
		a, n, err := s.FindUsers(ctx, fwt.UserFilter{Email: &email})
		require.NoError(t, err)
		require.Equal(t, len(a), 1)
		require.Equal(t, a[0].ID, uint(1))
		require.Equal(t, a[0].Username, "janedoe")
		require.Equal(t, a[0].Email, email)
		require.Equal(t, n, 1)
	})
}

func MustCreateUser(tb testing.TB, ctx context.Context, db *postgres.DB, user *fwt.User) *fwt.User {
	tb.Helper()
	err := postgres.NewUserService(db).CreateUser(ctx, user)
	require.NoError(tb, err)
	return user
}
