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
	return nil
}

func (s *UserService) UpdateUser(ctx context.Context, id uint, upd fwt.UserUpdate) (*fwt.User, error) {
	return nil, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	return nil
}
