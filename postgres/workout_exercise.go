package postgres

import (
	"context"
	"fmt"

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
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	we, err := findWorkoutExerciseByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return we, nil
}

func (s *WorkoutExerciseService) FindWorkoutExercises(ctx context.Context, filter fwt.WorkoutExerciseFilter) ([]*fwt.WorkoutExercise, int, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	return findWorkoutExercises(ctx, tx, filter)
}

func (s *WorkoutExerciseService) CreateWorkoutExercise(ctx context.Context, workoutExercise *fwt.WorkoutExercise) error {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := createWorkoutExercise(ctx, tx, workoutExercise); err != nil {
		return err
	}

	// Find the WEStatus with current we.ID to make sure to create
	// multiple weStatuses

	if err := createWEStatus(ctx, tx, &fwt.WEStatus{
		WorkoutExerciseID: workoutExercise.ID,
		Status:            "pending",
	}); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *WorkoutExerciseService) UpdateWorkoutExercise(ctx context.Context, id uint, upd fwt.WorkoutExerciseUpdate) (*fwt.WorkoutExercise, error) {
	return nil, nil
}

func (s *WorkoutExerciseService) DeleteWorkoutExercise(ctx context.Context, id uint) error {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := deleteWorkoutExercise(ctx, tx, id); err != nil {
		return err
	}

	return tx.Commit()
}

func createWorkoutExercise(ctx context.Context, tx *Tx, workoutExercise *fwt.WorkoutExercise) error {
	workoutExercise.CreatedAt = tx.now
	workoutExercise.UpdatedAt = workoutExercise.CreatedAt

	if err := workoutExercise.Validate(); err != nil {
		return err
	}

	/*
		workout, err := findWorkoutByID(ctx, tx, workoutExercise.WorkoutID)
		if err != nil {
			return err
			} */ // Hush-Hush

	exercise, err := findExerciseByID(ctx, tx, workoutExercise.ExerciseID)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO workout_exercise (workout_id, exercise_id, "order", created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5) RETURNING id
	`
	args := []interface{}{
		workoutExercise.WorkoutID,
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

func findWorkoutExercises(ctx context.Context, tx *Tx, filter fwt.WorkoutExerciseFilter) (_ []*fwt.WorkoutExercise, n int, err error) {
	where, args := []string{}, []interface{}{}
	argPos := 0

	if v := filter.ID; v != nil {
		argPos++
		where, args = append(where, fmt.Sprintf("id = $%d", argPos)), append(args, *v)
	}
	if v := filter.WorkoutID; v != nil {
		argPos++
		where, args = append(where, fmt.Sprintf("workout_id = $%d", argPos)), append(args, *v)
	}
	if v := filter.ExerciseID; v != nil {
		argPos++
		where, args = append(where, fmt.Sprintf("exercise_id = $%d", argPos)), append(args, *v)
	}
	if v := filter.Order; v != nil {
		argPos++
		where, args = append(where, fmt.Sprintf(`"order" = $%d`, argPos)), append(args, *v)
	}

	query := `
	SELECT id, workout_id, exercise_id, "order", created_at, updated_at, COUNT(*) OVER()
	FROM workout_exercise` + formatWhereClause(where) + ` ORDER BY id ASC` + formatLimitOffset(filter.Limit, filter.Offset)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	workoutExercises := make([]*fwt.WorkoutExercise, 0)
	for rows.Next() {
		var workoutExercise fwt.WorkoutExercise
		if err := rows.Scan(
			&workoutExercise.ID,
			&workoutExercise.WorkoutID,
			&workoutExercise.ExerciseID,
			&workoutExercise.Order,
			(*NullTime)(&workoutExercise.CreatedAt),
			(*NullTime)(&workoutExercise.UpdatedAt),
			&n,
		); err != nil {
			return nil, n, err
		}

		workoutExercises = append(workoutExercises, &workoutExercise)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return workoutExercises, n, nil
}

func findWorkoutExerciseByID(ctx context.Context, tx *Tx, id uint) (*fwt.WorkoutExercise, error) {
	a, _, err := findWorkoutExercises(ctx, tx, fwt.WorkoutExerciseFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, fwt.Errorf(fwt.ENOTFOUND, "Workout Exercise not found.")
	}

	return a[0], nil
}

func deleteWorkoutExercise(ctx context.Context, tx *Tx, id uint) error {
	if _, err := findWorkoutExerciseByID(ctx, tx, id); err != nil {
		return err
	}

	query := `
	DELETE FROM workout_exercise WHERE id = $1
	`
	if _, err := tx.ExecContext(ctx, query, id); err != nil {
		return err
	}

	return nil
}
