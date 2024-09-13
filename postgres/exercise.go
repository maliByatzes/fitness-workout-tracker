package postgres

import (
	"context"
	"fmt"

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
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	exercise, err := findExerciseByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return exercise, nil
}

func (s *ExerciseService) FindExerciseByName(ctx context.Context, name string) (*fwt.Exercise, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	exercise, err := findExerciseByName(ctx, tx, name)
	if err != nil {
		return nil, err
	}

	return exercise, nil
}

func (s *ExerciseService) FindExercises(ctx context.Context, filter fwt.ExerciseFilter) ([]*fwt.Exercise, int, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	return findExercises(ctx, tx, filter)
}

func (s *ExerciseService) CreateExercise(ctx context.Context, exercise *fwt.Exercise) error {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := createExercise(ctx, tx, exercise); err != nil {
		return err
	}

	return tx.Commit()
}

func createExercise(ctx context.Context, tx *Tx, exercise *fwt.Exercise) error {
	exercise.CreatedAt = tx.now
	exercise.UpdatedAt = exercise.CreatedAt

	if err := exercise.Validate(); err != nil {
		return err
	}

	args := []interface{}{
		exercise.Name,
		exercise.Description,
		(*NullTime)(&exercise.CreatedAt),
		(*NullTime)(&exercise.UpdatedAt),
	}
	query := `
	INSERT INTO exercise (name, description, created_at, updated_at)
	VALUES ($1, $2, $3, $4) RETURNING id
	`

	err := tx.QueryRowxContext(ctx, query, args...).Scan(&exercise.ID)
	if err != nil {
		return err
	}

	return nil
}

func findExerciseByID(ctx context.Context, tx *Tx, id uint) (*fwt.Exercise, error) {
	a, _, err := findExercises(ctx, tx, fwt.ExerciseFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &fwt.Error{Code: fwt.ENOTFOUND, Message: "Exercise not found."}
	}

	return a[0], nil
}

func findExerciseByName(ctx context.Context, tx *Tx, name string) (*fwt.Exercise, error) {
	a, _, err := findExercises(ctx, tx, fwt.ExerciseFilter{Name: &name})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &fwt.Error{Code: fwt.ENOTFOUND, Message: "Exercise not found."}
	}
	return a[0], nil
}

func findExercises(ctx context.Context, tx *Tx, filter fwt.ExerciseFilter) (_ []*fwt.Exercise, n int, err error) {
	where, args := []string{}, []interface{}{}
	argPos := 0

	if v := filter.ID; v != nil {
		argPos++
		where, args = append(where, fmt.Sprintf("id = $%d", argPos)), append(args, *v)
	}
	if v := filter.Name; v != nil {
		argPos++
		where, args = append(where, fmt.Sprintf("name = $%d", argPos)), append(args, *v)
	}

	query := `
	SELECT id, name, description, created_at, updated_at, COUNT(*) OVER()
	FROM exercise` + formatWhereClause(where) + ` ORDER BY id ASC` + formatLimitOffset(filter.Limit, filter.Limit)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	exercises := make([]*fwt.Exercise, 0)
	for rows.Next() {
		var exercise fwt.Exercise
		if err := rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.Description,
			(*NullTime)(&exercise.CreatedAt),
			(*NullTime)(&exercise.UpdatedAt),
			&n,
		); err != nil {
			return nil, n, err
		}

		exercises = append(exercises, &exercise)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return exercises, n, nil
}
