package postgres

import (
	"context"

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
	return nil, nil
}

func (s *WEStatusService) FindWEStatuses(ctx context.Context, filter fwt.WEStatusFilter) ([]*fwt.WEStatus, int, error) {
	return nil, 0, nil
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
