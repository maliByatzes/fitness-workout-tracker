package postgres

import (
	"context"

	"github.com/maliByatzes/fwt"
)

var _ fwt.WorkoutExerciseService = (*WorkoutExerciseService)(nil)

type WorkoutExerciseService struct {
	db *DB
}

func NewWorkoutExerciseService(db *DB) *WorkoutExerciseService {
	return &WorkoutExerciseService{db: db}
}

func (s *WorkoutExerciseService) FindWorkoutExerciseByID(ctx context.Context, id uint) (*fwt.WorkoutExercise, error) {
	return nil, nil
}

func (s *WorkoutExerciseService) FindWorkoutExercises(ctx context.Context, filter fwt.WorkoutExerciseFilter) ([]*fwt.WorkoutExercise, int, error) {
	return nil, 0, nil
}

func (s *WorkoutExerciseService) CreateWorkoutExercise(ctx context.Context, workout *fwt.WorkoutExercise) error {
	return nil
}

func (s *WorkoutExerciseService) UpdateWorkoutExercise(ctx context.Context, id uint, upd fwt.WorkoutExerciseUpdate) (*fwt.WorkoutExercise, error) {
	return nil, nil
}

func (s *WorkoutExerciseService) DeleteWorkoutExercise(ctx context.Context, id uint) error {
	return nil
}
