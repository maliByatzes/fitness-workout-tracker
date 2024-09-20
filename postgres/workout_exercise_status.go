package postgres

import (
	"context"
	"fmt"

	"github.com/maliByatzes/fwt"
)

var _ fwt.WEStatusService = (*WEStatusService)(nil)

type WEStatusService struct {
	db *DB
}

func NewWEStatusService(db *DB) *WEStatusService {
	return &WEStatusService{db: db}
}

func (s *WEStatusService) FindWEStatusByID(ctx context.Context, id uint) (*fwt.WEStatus, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	we, err := findWEStatusByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return we, nil
}

func (s *WEStatusService) FindWEStatusByWEID(ctx context.Context, id uint) (*fwt.WEStatus, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	we, err := findWEStatusByWEID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return we, nil
}

func (s *WEStatusService) FindWEStatuses(ctx context.Context, filter fwt.WEStatusFilter) ([]*fwt.WEStatus, int, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	return findWEStatuses(ctx, tx, filter)
}

func (s *WEStatusService) CreateWEStatus(ctx context.Context, we *fwt.WEStatus) error {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := createWEStatus(ctx, tx, we); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *WEStatusService) UpdateWEStatus(ctx context.Context, id uint, upd fwt.WEStatusUpdate) (*fwt.WEStatus, error) {
	return nil, nil
}

func (s *WEStatusService) DeleteWEStatus(ctx context.Context, id uint) error {
	return nil
}

func createWEStatus(ctx context.Context, tx *Tx, we *fwt.WEStatus) error {
	user := fwt.UserFromContext(ctx)
	if user == nil {
		return fwt.Errorf(fwt.ENOTAUTHORIZED, "You must logged in to create westatus.")
	}

	we.CreatedAt = tx.now
	we.UpdatedAt = we.CreatedAt

	if err := we.Validate(); err != nil {
		return err
	}

	query := `
	INSERT INTO workout_exercise_status (workout_exercise_id, status, comments, completed_at, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`
	args := []interface{}{
		we.WorkoutExerciseID,
		we.Status,
		we.Comments,
		(*NullTime)(&we.CompletedAt),
		(*NullTime)(&we.CreatedAt),
		(*NullTime)(&we.UpdatedAt),
	}

	err := tx.QueryRowxContext(ctx, query, args...).Scan(&we.ID)
	if err != nil {
		return err
	}

	return nil
}

func findWEStatuses(ctx context.Context, tx *Tx, filter fwt.WEStatusFilter) (_ []*fwt.WEStatus, n int, err error) {
	where, args := []string{}, []interface{}{}
	argPos := 0

	if v := filter.ID; v != nil {
		argPos++
		where, args = append(where, fmt.Sprintf("id = $%d", argPos)), append(args, *v)
	}
	if v := filter.WorkoutExerciseID; v != nil {
		argPos++
		where, args = append(where, fmt.Sprintf("workout_exercise_id = $%d", argPos)), append(args, *v)
	}
	if v := filter.Status; v != nil {
		argPos++
		where, args = append(where, fmt.Sprintf("status = $%d", argPos)), append(args, *v)
	}

	query := `
	SELECT id, workout_exercise_id, status, comments, completed_at, created_at, updated_at, COUNT(*) OVER()
	FROM workout_exercise_status` + formatWhereClause(where) + ` ORDER BY id ASC` + formatLimitOffset(filter.Limit, filter.Offset)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	weStatuses := make([]*fwt.WEStatus, 0)
	for rows.Next() {
		var weStatus fwt.WEStatus
		if err := rows.Scan(
			&weStatus.ID,
			&weStatus.WorkoutExerciseID,
			&weStatus.Status,
			&weStatus.Comments,
			(*NullTime)(&weStatus.CompletedAt),
			(*NullTime)(&weStatus.CreatedAt),
			(*NullTime)(&weStatus.UpdatedAt),
			&n,
		); err != nil {
			return nil, n, err
		}

		weStatuses = append(weStatuses, &weStatus)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return weStatuses, n, nil
}

func findWEStatusByID(ctx context.Context, tx *Tx, id uint) (*fwt.WEStatus, error) {
	a, _, err := findWEStatuses(ctx, tx, fwt.WEStatusFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, fwt.Errorf(fwt.ENOTFOUND, "WEStatus not found.")
	}

	return a[0], nil
}

func findWEStatusByWEID(ctx context.Context, tx *Tx, id uint) (*fwt.WEStatus, error) {
	a, _, err := findWEStatuses(ctx, tx, fwt.WEStatusFilter{WorkoutExerciseID: &id})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, fwt.Errorf(fwt.ENOTFOUND, "WEStatsu not found.")
	}

	return a[0], nil
}
