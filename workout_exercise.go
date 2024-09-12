package fwt

import (
	"context"
	"time"
)

type WorkoutExercise struct {
	ID         uint      `json:"id"`
	WorkoutID  uint      `json:"workout_id"`
	ExerciseID uint      `json:"exercise_id"`
	Order      uint      `json:"order"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (we *WorkoutExercise) Validate() error {
	if we.WorkoutID <= 0 {
		return Errorf(EINVALID, "WorkoutID is required.")
	}

	if we.ExerciseID <= 0 {
		return Errorf(EINVALID, "ExerciseID is required.")
	}

	if we.Order <= 0 {
		return Errorf(EINVALID, "Order is required.")
	}

	return nil
}

type WorkoutExerciseService interface {
	FindWorkoutExerciseByID(context.Context, uint) (*WorkoutExercise, error)
	FindWorkoutExercises(context.Context, WorkoutExerciseFilter) ([]*WorkoutExercise, int, error)
	CreateWorkoutExercise(context.Context, *WorkoutExercise) error
	UpdateWorkfoutExercise(context.Context, uint, WorkoutExerciseUpdate) (*WorkoutExercise, error)
	DeleteWorkoutExercise(context.Context, uint) error
}

type WorkoutExerciseFilter struct {
	ID         *uint `json:"id"`
	WorkoutID  *uint `json:"workout_id"`
	ExerciseID *uint `json:"exercise_id"`
	Order      *uint `json:"order"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type WorkoutExerciseUpdate struct {
	Order *uint `json:"order"`
}
