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
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := createWorkout(ctx, tx, workout); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *WorkoutService) UpdateWorkout(ctx context.Context, id uint, upd fwt.WorkoutUpdate) (*fwt.Workout, error) {
	return nil, nil
}

func (s *WorkoutService) DeleteWorkout(ctx context.Context, id uint) error {
	return nil
}

func createWorkout(ctx context.Context, tx *Tx, workout *fwt.Workout) error {
	workout.CreatedAt = tx.now
	workout.UpdatedAt = workout.CreatedAt

	if err := workout.Validate(); err != nil {
		return err
	}

	query := `
	INSERT INTO workout (user_id, name, scheduled_date, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5) RETURNING id
	`
	args := []interface{}{
		workout.UserID,
		workout.Name,
		workout.ScheduledDate,
		(*NullTime)(&workout.CreatedAt),
		(*NullTime)(&workout.UpdatedAt),
	}

	err := tx.QueryRowxContext(ctx, query, args...).Scan(&workout.ID)
	if err != nil {
		return err
	}

	return nil
}
