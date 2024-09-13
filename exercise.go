package fwt

import (
	"context"
	"time"
)

type Exercise struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (e *Exercise) Validate() error {
	if e.Name == "" {
		return Errorf(EINVALID, "Name is required.")
	}

	if e.Description == "" {
		return Errorf(EINVALID, "Description is required.")
	}

	return nil
}

type ExerciseService interface {
	FindExerciseByID(context.Context, uint) (*Exercise, error)
	FindExerciseByName(context.Context, string) (*Exercise, error)
	FindExercises(context.Context, ExerciseFilter) ([]*Exercise, int, error)
	CreateExercise(context.Context, *Exercise) error
}

type ExerciseFilter struct {
	ID   *uint   `json:"id"`
	Name *string `json:"name"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}
