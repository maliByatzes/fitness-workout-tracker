package postgres

import (
	"context"

	"github.com/maliByatzes/fwt"
)

var _ fwt.ProfileService = (*ProfileService)(nil)

type ProfileService struct {
	db *DB
}

func NewProfileService(db *DB) *ProfileService {
	return &ProfileService{db: db}
}

func (s *ProfileService) FindProfileByID(ctx context.Context, id uint) (*fwt.Profile, error) {
	return nil, nil
}

func (s *ProfileService) FindDials(ctx context.Context, filter fwt.ProfileFilter) ([]*fwt.Profile, int, error) {
	return nil, 0, nil
}

func (s *ProfileService) CreateProfile(ctx context.Context, profile *fwt.Profile) error {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := createProfile(ctx, tx, profile); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *ProfileService) UpdateProfile(ctx context.Context, id uint, upd fwt.ProfileUpdate) (*fwt.Profile, error) {
	return nil, nil
}

func (s *ProfileService) DeleteProfile(ctx context.Context, id uint) error {
	return nil
}

func createProfile(ctx context.Context, tx *Tx, profile *fwt.Profile) error {
	profile.CreatedAt = tx.now
	profile.UpdatedAt = profile.CreatedAt

	if err := profile.Validate(); err != nil {
		return err
	}

	query := `
	INSERT INTO profile (user_id, first_name, last_name, date_of_birth, gender, height, weight, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id
	`
	args := []interface{}{
		profile.UserID,
		profile.FirstName,
		profile.LastName,
		profile.DateOfBirth,
		profile.Gender,
		profile.Height,
		profile.Weight,
		(*NullTime)(&profile.CreatedAt),
		(*NullTime)(&profile.UpdatedAt),
	}

	err := tx.QueryRowxContext(ctx, query, args...).Scan(&profile.ID)
	if err != nil {
		return err
	}

	return nil
}
