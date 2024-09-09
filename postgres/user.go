package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/maliByatzes/fwt"
)

var _ fwt.UserService = (*UserService)(nil)

type UserService struct {
	db *DB
}

func NewUserService(db *DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) FindUserByID(ctx context.Context, id uint) (*fwt.User, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	user, err := findUserByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Authenticate(ctx context.Context, username, password string) (*fwt.User, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	user, err := findUserByUsername(ctx, tx, username)
	if err != nil {
		return nil, err
	}

	if err := user.VerifyPassword(password, user.HashedPassword); err != nil {
		return nil, &fwt.Error{Code: fwt.ENOTAUTHORIZED, Message: "Incorrect credentials"}
	}

	return user, nil
}

func (s *UserService) FindUsers(ctx context.Context, filter fwt.UserFilter) ([]*fwt.User, int, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	return findUsers(ctx, tx, filter)
}

func (s *UserService) CreateUser(ctx context.Context, user *fwt.User) error {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := createUser(ctx, tx, user); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *UserService) UpdateUser(ctx context.Context, id uint, upd fwt.UserUpdate) (*fwt.User, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	user, err := updateUser(ctx, tx, id, upd)
	if err != nil {
		return user, err
	} else if err := tx.Commit(); err != nil {
		return user, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := deleteUser(ctx, tx, id); err != nil {
		return err
	}

	return tx.Commit()
}

func createUser(ctx context.Context, tx *Tx, user *fwt.User) error {
	user.CreatedAt = tx.now
	user.UpdatedAt = user.CreatedAt

	if err := user.Validate(); err != nil {
		return err
	}

	query := `
	INSERT INTO "user" (username, email, hashed_password, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5) RETURNING id
	`
	args := []interface{}{user.Username, user.Email, user.HashedPassword, (*NullTime)(&user.CreatedAt), (*NullTime)(&user.UpdatedAt)}

	err := tx.QueryRowxContext(ctx, query, args...).Scan(&user.ID)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "user_username_key"`:
			return fwt.Errorf(fwt.ECONFLICT, "This username already exists.")
		case err.Error() == `pq: duplicate key value violates unique constraint "user_email_key"`:
			return fwt.Errorf(fwt.ECONFLICT, "This email already exists.")
		default:
			return err
		}
	}

	return nil
}

func findUserByID(ctx context.Context, tx *Tx, id uint) (*fwt.User, error) {
	a, _, err := findUsers(ctx, tx, fwt.UserFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &fwt.Error{Code: fwt.ENOTFOUND, Message: "User not found."}
	}
	return a[0], nil
}

func findUserByUsername(ctx context.Context, tx *Tx, username string) (*fwt.User, error) {
	a, _, err := findUsers(ctx, tx, fwt.UserFilter{Username: &username})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &fwt.Error{Code: fwt.ENOTFOUND, Message: "User not found."}
	}
	return a[0], nil
}

func findUserByEmail(ctx context.Context, tx *Tx, email string) (*fwt.User, error) {
	a, _, err := findUsers(ctx, tx, fwt.UserFilter{Email: &email})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &fwt.Error{Code: fwt.ENOTFOUND, Message: "User not found."}
	}
	return a[0], nil
}

func findUsers(ctx context.Context, tx *Tx, filter fwt.UserFilter) (_ []*fwt.User, n int, err error) {
	where, args := []string{}, []interface{}{}
	argPosition := 0

	if v := filter.ID; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("id = $%d", argPosition)), append(args, *v)
	}

	if v := filter.Username; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("username = $%d", argPosition)), append(args, *v)
	}

	if v := filter.Email; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("email = $%d", argPosition)), append(args, *v)
	}

	query := `SELECT id, username, email, hashed_password, created_at, updated_at, COUNT(*) OVER() FROM "user"` + formatWhereClause(where) +
		` ORDER BY id ASC` + formatLimitOffset(filter.Limit, filter.Offset)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	users := make([]*fwt.User, 0)
	for rows.Next() {
		var user fwt.User
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.HashedPassword,
			(*NullTime)(&user.CreatedAt),
			(*NullTime)(&user.UpdatedAt),
			&n,
		); err != nil {
			return nil, n, err
		}

		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, n, nil
}

func updateUser(ctx context.Context, tx *Tx, id uint, upd fwt.UserUpdate) (*fwt.User, error) {
	user, err := findUserByID(ctx, tx, id)
	if err != nil {
		return user, err
	}

	if v := upd.Username; v != nil {
		user.Username = *v
	}

	if v := upd.Email; v != nil {
		user.Email = *v
	}

	user.UpdatedAt = tx.now

	if err := user.Validate(); err != nil {
		return user, err
	}

	args := []interface{}{
		user.Username,
		user.Email,
		user.ID,
	}
	query := `
	UPDATE "user" SET username = $1, email = $2
	WHERE ID = $3
	`

	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return user, err
	}

	return user, nil
}

func deleteUser(ctx context.Context, tx *Tx, id uint) error {
	if _, err := findUserByID(ctx, tx, id); err != nil {
		return err
	}

	query := `
	DELETE FROM "user" WHERE id = $1
	`
	if _, err := tx.ExecContext(ctx, query, id); err != nil {
		return err
	}

	return nil
}

func formatLimitOffset(limit, offset int) string {
	if limit > 0 && offset > 0 {
		return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
	} else if limit > 0 {
		return fmt.Sprintf("LIMIT %d", limit)
	} else if offset > 0 {
		return fmt.Sprintf("OFFSET %d", offset)
	}
	return ""
}

func formatWhereClause(where []string) string {
	if len(where) == 0 {
		return ""
	}
	return " WHERE " + strings.Join(where, " AND ")
}
