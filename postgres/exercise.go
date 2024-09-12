package postgres

import (
	"context"

	"github.com/maliByatzes/fwt"
)

var _ fwt.ExerciseService = (*ExerciseService)(nil)

type ExerciseService struct {
	db *DB
}

func NewExerciseService(db *DB) *ExerciseService {
	return &ExerciseService{db: db}
}

func (s *ExerciseService) FindExerciseByID(ctx context.Context, id uint) (*fwt.Exercise, error) {
	return nil, nil
}

func (s *ExerciseService) FindExerciseByName(ctx context.Context, name string) (*fwt.Exercise, error) {
	return nil, nil
}

func (s *ExerciseService) FindExercises(ctx context.Context, filter fwt.ExerciseFilter) ([]*fwt.Exercise, int, error) {
	return nil, 0, nil
}
