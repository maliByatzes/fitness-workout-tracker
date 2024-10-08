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

func (s *WorkoutService) CreateWorkout(ctx context.Context, workout *fwt.Workout) error {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := createWorkout(ctx, tx, workout); err != nil {
		return err
	}

	for _, ex := range workout.Exercises {
		exercise, err := findExerciseByName(ctx, tx, ex.Name)
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
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	workout, err := updateWorkout(ctx, tx, id, upd)
	if err != nil {
		return workout, err
	} else if err := tx.Commit(); err != nil {
		return workout, err
	}

	return workout, nil
}

func (s *WorkoutService) RemoveExercisesFromWorkout(ctx context.Context, id uint, exercises []string) (*fwt.Workout, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	workout, err := findWorkoutByID(ctx, tx, id)
	if err != nil {
		return nil, err
	} else if workout.UserID != fwt.UserIDFromContext(ctx) {
		return nil, fwt.Errorf(fwt.ENOTAUTHORIZED, "You are not allowed to modify this workout.")
	}

	for _, exName := range exercises {
		for index, e := range workout.Exercises {
			if e.Name == exName {
				err := removeExerciseFromWorkout(ctx, tx, workout, e)
				if err != nil {
					return workout, err
				}

				workout.Exercises = append(workout.Exercises[:index], workout.Exercises[index+1:]...)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return workout, err
	}

	return workout, nil
}

func (s *WorkoutService) AddExercisesToWorkout(ctx context.Context, id uint, exercises []string) (*fwt.Workout, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	workout, err := findWorkoutByID(ctx, tx, id)
	if err != nil {
		return nil, err
	} else if workout.UserID != fwt.UserIDFromContext(ctx) {
		return nil, fwt.Errorf(fwt.ENOTAUTHORIZED, "You are not allowed to modify this workout.")
	}

	for _, exName := range exercises {
		exercise, err := findExerciseByName(ctx, tx, exName)
		if err != nil {
			return workout, err
		}

		if !implContains2(workout.Exercises, exercise) {
			err = addExerciseToWorkout(ctx, tx, workout, exercise)
			if err != nil {
				return workout, nil
			}

			workout.Exercises = append(workout.Exercises, exercise)
		}
	}

	if err := tx.Commit(); err != nil {
		return workout, err
	}

	return workout, nil
}

func (s *WorkoutService) DeleteWorkout(ctx context.Context, id uint) error {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := deleteWorkout(ctx, tx, id); err != nil {
		return err
	}

	return tx.Commit()
}

func createWorkout(ctx context.Context, tx *Tx, workout *fwt.Workout) error {
	user := fwt.UserFromContext(ctx)
	if user == nil {
		return fwt.Errorf(fwt.ENOTAUTHORIZED, "You must be logged to create workout")
	}
	workout.UserID = fwt.UserIDFromContext(ctx)

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
		where, args = append(where, fmt.Sprintf("w.id = $%d", argPos)), append(args, *v)
	}
	if v := filter.UserID; v != nil {
		argPos++
		where, args = append(where, fmt.Sprintf("w.user_id = $%d", argPos)), append(args, *v)
	}
	if v := filter.Name; v != nil {
		argPos++
		where, args = append(where, fmt.Sprintf("w.name = $%d", argPos)), append(args, *v)
	}
	if v := filter.ScheduledDate; v != nil {
		argPos++
		where, args = append(where, fmt.Sprintf("w.scheduled_date = $%d", argPos)), append(args, *v)
	}

	query := `
	SELECT w.id, w.user_id, w.name, w.scheduled_date, w.created_at, w.updated_at, e.id, e.name, e.description, e.created_at, e.updated_at, COUNT(*) OVER()
	FROM workout AS w
	INNER JOIN workout_exercise AS we ON we.workout_id = w.id
	INNER JOIN exercise as e ON e.id = we.exercise_id` + formatWhereClause(where) + ` ORDER BY w.id ASC` + formatLimitOffset(filter.Limit, filter.Offset)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	workouts := make([]*fwt.Workout, 0)
	for rows.Next() {
		var workout fwt.Workout
		var exercise fwt.Exercise
		if err := rows.Scan(
			&workout.ID,
			&workout.UserID,
			&workout.Name,
			&workout.ScheduledDate,
			(*NullTime)(&workout.CreatedAt),
			(*NullTime)(&workout.UpdatedAt),
			&exercise.ID,
			&exercise.Name,
			&exercise.Description,
			(*NullTime)(&exercise.CreatedAt),
			(*NullTime)(&exercise.UpdatedAt),
			&n,
		); err != nil {
			return nil, n, err
		}

		if index := implContains(workouts, &workout); index != -1 {
			workouts[index].Exercises = append(workouts[index].Exercises, &exercise)
		} else {
			workout.Exercises = append(workout.Exercises, &exercise)
			workouts = append(workouts, &workout)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return workouts, n, nil
}

func updateWorkout(ctx context.Context, tx *Tx, id uint, upd fwt.WorkoutUpdate) (*fwt.Workout, error) {
	workout, err := findWorkoutByID(ctx, tx, id)
	if err != nil {
		return workout, err
	} else if workout.UserID != fwt.UserIDFromContext(ctx) {
		return nil, fwt.Errorf(fwt.ENOTAUTHORIZED, "You are not allowed to update this workout.")
	}

	if v := upd.Name; v != nil {
		workout.Name = *v
	}
	if v := upd.ScheduledDate; v != nil {
		workout.ScheduledDate = *v
	}
	workout.UpdatedAt = tx.now

	if err := workout.Validate(); err != nil {
		return workout, err
	}

	args := []interface{}{
		workout.Name,
		workout.ScheduledDate,
		workout.UpdatedAt,
		workout.ID,
		workout.UserID,
	}
	query := `
	UPDATE workout SET name = $1, scheduled_date = $2, updated_at = $3
	WHERE id = $4 AND user_id = $5
	`

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return workout, err
	}

	return workout, nil
}

func removeExerciseFromWorkout(ctx context.Context, tx *Tx, workout *fwt.Workout, exercise *fwt.Exercise) error {
	a, _, err := findWorkoutExercises(ctx, tx, fwt.WorkoutExerciseFilter{WorkoutID: &workout.ID, ExerciseID: &exercise.ID})
	if err != nil {
		return err
	} else if len(a) == 0 {
		return fwt.Errorf(fwt.ENOTFOUND, "Workout_Exercise not found.")
	}

	err = deleteWorkoutExercise(ctx, tx, a[0].ID)
	if err != nil {
		return err
	}

	return nil
}

func addExerciseToWorkout(ctx context.Context, tx *Tx, wokout *fwt.Workout, exercise *fwt.Exercise) error {
	if err := createWorkoutExercise(ctx, tx, &fwt.WorkoutExercise{
		WorkoutID:  wokout.ID,
		ExerciseID: exercise.ID,
		Order:      1, // Hard-code for now ...
	}); err != nil {
		return err
	}

	return nil
}

func deleteWorkout(ctx context.Context, tx *Tx, id uint) error {
	workout, err := findWorkoutByID(ctx, tx, id)
	if err != nil {
		return err
	} else if workout.UserID != fwt.UserIDFromContext(ctx) {
		return fwt.Errorf(fwt.ENOTAUTHORIZED, "You are not allowed to delete this workout.")
	}

	workoutExercises, _, err := findWorkoutExercises(ctx, tx, fwt.WorkoutExerciseFilter{WorkoutID: &workout.ID})
	if err != nil {
		return err
	}

	for _, we := range workoutExercises {
		err := deleteWorkoutExercise(ctx, tx, we.ID)
		if err != nil {
			return err
		}
	}

	args := []interface{}{
		workout.ID,
		workout.UserID,
	}
	query := `
	DELETE FROM workout WHERE id = $1 AND user_id = $2
	`

	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func implContains(ws []*fwt.Workout, w *fwt.Workout) int {
	for index, value := range ws {
		if value.ID == w.ID {
			return index
		}
	}
	return -1
}

func implContains2(exs []*fwt.Exercise, ex *fwt.Exercise) bool {
	for _, value := range exs {
		if value.Name == ex.Name {
			return true
		}
	}

	return false
}
