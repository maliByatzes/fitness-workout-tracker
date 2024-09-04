package postgres

import (
	"context"

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
	return nil, nil
}

func (s *UserService) FindUsers(ctx context.Context, filter fwt.UserFilter) ([]*fwt.User, int, error) {
	return nil, 0, nil
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
	return nil, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	return nil
}

func createUser(ctx context.Context, tx *Tx, user *fwt.User) error {
	user.CreatedAt = tx.now
	user.UpdatedAt = user.CreatedAt

	if err := user.Validate(); err != nil {
		return err
	}

	query := `
	INSERT INTO user (username, email, hashed_password, created_at, updated_at)
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
