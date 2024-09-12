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

func (s *WorkoutExerciseService) CreateWorkoutExercise(ctx context.Context, workoutExercise *fwt.WorkoutExercise) error {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := createWorkoutExercise(ctx, tx, workoutExercise); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *WorkoutExerciseService) UpdateWorkoutExercise(ctx context.Context, id uint, upd fwt.WorkoutExerciseUpdate) (*fwt.WorkoutExercise, error) {
	return nil, nil
}

func (s *WorkoutExerciseService) DeleteWorkoutExercise(ctx context.Context, id uint) error {
	return nil
}

func createWorkoutExercise(ctx context.Context, tx *Tx, workoutExercise *fwt.WorkoutExercise) error {
	workoutExercise.CreatedAt = tx.now
	workoutExercise.UpdatedAt = workoutExercise.CreatedAt

	if err := workoutExercise.Validate(); err != nil {
		return err
	}

	workout, err := findWorkoutByID(ctx, tx, workoutExercise.WorkoutID)
	if err != nil {
		return err
	}

	exercise, err := findExerciseByID(ctx, tx, workoutExercise.ExerciseID)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO workout_exercise (workout_id, exercise_id, "order", created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5) RETURNING id
	`
	args := []interface{}{
		workout.ID,
		exercise.ID,
		workoutExercise.Order,
		(*NullTime)(&workoutExercise.CreatedAt),
		(*NullTime)(&workoutExercise.UpdatedAt),
	}

	err = tx.QueryRowxContext(ctx, query, args...).Scan(&workoutExercise.ID)
	if err != nil {
		return err
	}

	return nil
}
