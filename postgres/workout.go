package postgres

import (
	"context"

	"github.com/maliByatzes/fwt"
)

var _ fwt.WorkoutService = (*WorkoutService)(nil)

type WorkoutService struct {
	db *DB
}

func NewWorkoutService(db *DB) *WorkoutService {
	return &WorkoutService{db: db}
}

func (s *WorkoutService) FindWorkoutByID(ctx context.Context, id uint) (*fwt.Workout, error) {
	return nil, nil
}

func (s *WorkoutService) FindWorkouts(ctx context.Context, filter fwt.WorkoutFilter) ([]*fwt.Workout, int, error) {
	return nil, 0, nil
}

func (s *WorkoutService) CreateWorkout(ctx context.Context, workout *fwt.Workout) error {
	return nil
}

func (s *WorkoutService) UpdateWorkout(ctx context.Context, id uint, upd fwt.WorkoutUpdate) (*fwt.Workout, error) {
	return nil, nil
}

func (s *WorkoutService) DeleteWorkout(ctx context.Context, id uint) error {
	return nil
}
