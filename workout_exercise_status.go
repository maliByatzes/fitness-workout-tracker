package fwt

import (
	"context"
	"time"
)

type WEStatus struct {
	ID                uint      `json:"id"`
	WorkoutExerciseID uint      `json:"workout_exercise_id"`
	Status            string    `json:"status"` // `binding:"oneof=pending,completed"`
	Comments          string    `json:"comments"`
	CompletedAt       time.Time `json:"completed_at"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (wes *WEStatus) Validate() error {
	if wes.WorkoutExerciseID <= 0 {
		return Errorf(EINVALID, "WorkoutExerciseID is required.")
	}

	if wes.Status == "" {
		return Errorf(EINVALID, "Status is required.")
	}

	return nil
}

type WEStatusService interface {
	FindWEStatusByID(context.Context, uint) (*WEStatus, error)
	FindWEStatusByWEID(context.Context, uint) (*WEStatus, error)
	FindWEStatuses(context.Context, WEStatusFilter) ([]*WEStatus, int, error)
	CreateWEStatus(context.Context, *WEStatus) error
	UpdateWEStatus(context.Context, uint, WEStatusUpdate) (*WEStatus, error)
	DeleteWEStatus(context.Context, uint) error
}

type WEStatusFilter struct {
	ID                *uint   `json:"id"`
	WorkoutExerciseID *uint   `json:"workout_exercise_id"`
	Status            *string `json:"status"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type WEStatusUpdate struct {
	Status      *string    `json:"status"`
	Comments    *string    `json:"comments"`
	CompletedAt *time.Time `json:"completed_at"`
}
