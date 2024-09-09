package mock

import (
	"context"

	"github.com/maliByatzes/fwt"
)

var _ fwt.UserService = (*UserService)(nil)

type UserService struct {
	FindUserbyIDFn func(ctx context.Context, id uint) (*fwt.User, error)
	AuthenticateFn func(ctx context.Context, username, password string) (*fwt.User, error)
	FindUsersFn    func(ctx context.Context, filter fwt.UserFilter) ([]*fwt.User, int, error)
	CreateUserFn   func(ctx context.Context, user *fwt.User) error
	UpdateUserFn   func(ctx context.Context, id uint, upd fwt.UserUpdate) (*fwt.User, error)
	DeleteUserFn   func(ctx context.Context, id uint) error
}

func (s *UserService) FindUserByID(ctx context.Context, id uint) (*fwt.User, error) {
	return s.FindUserbyIDFn(ctx, id)
}

func (s *UserService) Authenticate(ctx context.Context, username, password string) (*fwt.User, error) {
	return s.AuthenticateFn(ctx, username, password)
}

func (s *UserService) FindUsers(ctx context.Context, filter fwt.UserFilter) ([]*fwt.User, int, error) {
	return s.FindUsersFn(ctx, filter)
}

func (s *UserService) CreateUser(ctx context.Context, user *fwt.User) error {
	return s.CreateUserFn(ctx, user)
}

func (s *UserService) UpdateUser(ctx context.Context, id uint, upd fwt.UserUpdate) (*fwt.User, error) {
	return s.UpdateUserFn(ctx, id, upd)
}

func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	return s.DeleteUserFn(ctx, id)
}
