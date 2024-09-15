package postgres

import (
	"context"
	"fmt"

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
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	workout, err := findWorkoutByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return workout, nil
}

func (s *WorkoutService) FindWorkoutByIDUserID(ctx context.Context, id uint, userID uint) (*fwt.Workout, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	workout, err := findWorkoutByIDUserID(ctx, tx, id, userID)
	if err != nil {
		return nil, err
	}

	return workout, nil
}

func (s *WorkoutService) FindWorkouts(ctx context.Context, filter fwt.WorkoutFilter) ([]*fwt.Workout, int, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	return findWorkouts(ctx, tx, filter)
}

func (s *WorkoutService) CreateWorkout(ctx context.Context, workout *fwt.Workout, exercises []string) error {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := createWorkout(ctx, tx, workout); err != nil {
		return err
	}

	for _, exName := range exercises {
		exercise, err := findExerciseByName(ctx, tx, exName)
		if err != nil {
			return err
		}

		if err := createWorkoutExercise(ctx, tx, &fwt.WorkoutExercise{
			WorkoutID:  workout.ID,
			ExerciseID: exercise.ID,
			Order:      1, // Hard-code for now...
		}); err != nil {
			return err
		}
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

func findWorkoutByID(ctx context.Context, tx *Tx, id uint) (*fwt.Workout, error) {
	a, _, err := findWorkouts(ctx, tx, fwt.WorkoutFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &fwt.Error{Code: fwt.ENOTFOUND, Message: "Workout not found."}
	}

	return a[0], nil
}

func findWorkoutByIDUserID(ctx context.Context, tx *Tx, id uint, userID uint) (*fwt.Workout, error) {
	a, _, err := findWorkouts(ctx, tx, fwt.WorkoutFilter{ID: &id, UserID: &userID})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &fwt.Error{Code: fwt.ENOTFOUND, Message: "Workout not found."}
	}

	return a[0], nil
}

func findWorkouts(ctx context.Context, tx *Tx, filter fwt.WorkoutFilter) (_ []*fwt.Workout, n int, err error) {
	where, args := []string{}, []interface{}{}
	argPos := 0

	if v := filter.ID; v != nil {
		argPos++
		where, args = append(where, fmt.Sprintf("id = $%d", argPos)), append(args, *v)
	}
	if v := filter.UserID; v != nil {
		argPos++
		where, args = append(where, fmt.Sprintf("user_id = $%d", argPos)), append(args, *v)
	}
	if v := filter.Name; v != nil {
		argPos++
		where, args = append(where, fmt.Sprintf("name = $%d", argPos)), append(args, *v)
	}
	if v := filter.ScheduledDate; v != nil {
		argPos++
		where, args = append(where, fmt.Sprintf("scheduled_date = $%d", argPos)), append(args, *v)
	}

	query := `
	SELECT id, user_id, name, scheduled_date, created_at, updated_at, COUNT(*) OVER()
	FROM workout` + formatWhereClause(where) + ` ORDER BY id ASC` + formatLimitOffset(filter.Limit, filter.Offset)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	workouts := make([]*fwt.Workout, 0)
	for rows.Next() {
		var workout fwt.Workout
		if err := rows.Scan(
			&workout.ID,
			&workout.UserID,
			&workout.Name,
			&workout.ScheduledDate,
			(*NullTime)(&workout.CreatedAt),
			(*NullTime)(&workout.UpdatedAt),
			&n,
		); err != nil {
			return nil, n, err
		}

		workouts = append(workouts, &workout)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return workouts, n, nil
}
