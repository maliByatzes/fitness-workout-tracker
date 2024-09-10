package fwt

import (
	"context"
	"time"
)

type Profile struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	DateOfBirth time.Time `json:"dob"`
	Gender      string    `json:"gender"`
	Height      float64   `json:"height"`
	Weight      float64   `json:"weight"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProfileService interface {
	FindProfileByID(ctx context.Context, id uint) (*Profile, error)
	FindDials(ctx context.Context, filter ProfileFilter) ([]*Profile, int, error)
	CreateProfile(ctx context.Context, profile *Profile) error
	UpdateProfile(ctx context.Context, id uint, upd ProfileUpdate) (*Profile, error)
	DeleteProfile(ctx context.Context, id uint) error
}

type ProfileFilter struct {
	ID          *uint      `json:"id"`
	FirstName   *string    `json:"first_name"`
	LastName    *string    `json:"last_name"`
	DateOfBirth *time.Time `json:"dob"`
	Gender      *string    `json:"gender"`
	Height      *float64   `json:"height"`
	Weight      *float64   `json:"weight"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type ProfileUpdate struct {
	FirstName   *string    `json:"first_name"`
	LastName    *string    `json:"last_name"`
	DateOfBirth *time.Time `json:"dob"`
	Gender      *string    `json:"gender"`
	Height      *float64   `json:"height"`
	Weight      *float64   `json:"weight"`
}
