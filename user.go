package fwt

import (
	"context"
	"time"
)

type User struct {
	ID             uint      `json:"id"`
	Username       string    `json:"username,omitempty"`
	Email          string    `json:"email,omitempty"`
	HashedPassword string    `json:"-" db:"hashed_password"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type Profile struct {
	UserID      uint      `json:"userID"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	DateOfBirth time.Time `json:"dob"`
	Gender      string    `json:"gender"`
	Height      float64   `json:"height"`
	Weight      float64   `json:"weight"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (u *User) Validate() error {
	if u.Username == "" {
		return Errorf(EINVALID, "Username is required.")
	}
	if u.Email == "" {
		return Errorf(EINVALID, "Email is required.")
	}
	return nil
}

type UserService interface {
	FindUserByID(ctx context.Context, id uint) (*User, error)
	FindUsers(ctx context.Context, filter UserFilter) ([]*User, int, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, id uint, upd UserUpdate) (*User, error)
	DeleteUser(ctx context.Context, id uint) error
}

type UserFilter struct {
	ID       *uint   `json:"id"`
	Username *string `json:"username"`
	Email    *string `json:"email"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type UserUpdate struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
}
