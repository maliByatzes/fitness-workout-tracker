package fwt

import (
	"context"
	"time"
)

type Workout struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"user_id"`
	Name          string    `json:"name"`
	ScheduledDate time.Time `json:"scheduled_date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (w *Workout) Validate() error {
	if w.UserID <= uint(0) {
		return Errorf(EINVALID, "UserID is required.")
	}

	if w.Name == "" {
		return Errorf(EINVALID, "Name is required.")
	}

	if w.ScheduledDate.IsZero() {
		return Errorf(EINVALID, "Scheduled Date is required.")
	}

	if !w.ScheduledDate.After(time.Now()) {
		return Errorf(EINVALID, "Scheduled Date is invalid.")
	}

	return nil
}

type WorkoutService interface {
	FindWorkoutByID(context.Context, uint) (*Workout, error)
	FindWorkouts(context.Context, WorkoutFilter) ([]*Workout, int, error)
	CreateWorkout(context.Context, *Workout, []string) error
	UpdateWorkout(context.Context, uint, WorkoutUpdate) (*Workout, error)
	DeleteWorkout(context.Context, uint) error
}

type WorkoutFilter struct {
	ID            *uint      `json:"id"`
	UserID        *uint      `json:"user_id"`
	Name          *string    `json:"name"`
	ScheduledDate *time.Time `json:"scheduled_date"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type WorkoutUpdate struct {
	Name          *string    `json:"name"`
	ScheduledDate *time.Time `json:"scheduled_date"`
}
