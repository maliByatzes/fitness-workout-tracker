package postgres

import (
	"context"
	"fmt"

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
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	profile, err := findProfileByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *ProfileService) FindProfiles(ctx context.Context, filter fwt.ProfileFilter) ([]*fwt.Profile, int, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	return findProfiles(ctx, tx, filter)
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

	exPr, err := findProfileByUserID(ctx, tx, profile.UserID)
	if err != nil && fwt.ErrorCode(err) != fwt.ENOTFOUND && fwt.ErrorMessage(err) != "Profile not found." {
		return err
	}
	if exPr != nil {
		return &fwt.Error{Code: fwt.ECONFLICT, Message: "Profile already exists."}
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

	err = tx.QueryRowxContext(ctx, query, args...).Scan(&profile.ID)
	if err != nil {
		return err
	}

	return nil
}

func findProfileByID(ctx context.Context, tx *Tx, id uint) (*fwt.Profile, error) {
	a, _, err := findProfiles(ctx, tx, fwt.ProfileFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &fwt.Error{Code: fwt.ENOTFOUND, Message: "Profile not found."}
	}
	return a[0], nil
}

func findProfileByUserID(ctx context.Context, tx *Tx, userID uint) (*fwt.Profile, error) {
	a, _, err := findProfiles(ctx, tx, fwt.ProfileFilter{UserID: &userID})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &fwt.Error{Code: fwt.ENOTFOUND, Message: "Profile not found."}
	}
	return a[0], nil
}

func findProfiles(ctx context.Context, tx *Tx, filter fwt.ProfileFilter) (_ []*fwt.Profile, n int, err error) {
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
	if v := filter.FirstName; v != nil {
		argPos++
		where, args = append(where, fmt.Sprintf("first_name = $%d", argPos)), append(args, *v)
	}
	if v := filter.LastName; v != nil {
		argPos++
		where, args = append(where, fmt.Sprintf("last_name = $%d", argPos)), append(args, *v)
	}
	if v := filter.DateOfBirth; v != nil {
		argPos++
		where, args = append(where, fmt.Sprintf("date_of_birth = $%d", argPos)), append(args, *v)
	}
	if v := filter.Gender; v != nil {
		argPos++
		where, args = append(where, fmt.Sprintf("gender = $%d", argPos)), append(args, *v)
	}
	if v := filter.Height; v != nil {
		argPos++
		where, args = append(where, fmt.Sprintf("height = $%d", argPos)), append(args, *v)
	}
	if v := filter.Weight; v != nil {
		argPos++
		where, args = append(where, fmt.Sprintf("weight = $%d", argPos)), append(args, *v)
	}

	query := `
	SELECT id, user_id, first_name, last_name, date_of_birth, gender, height, weight, created_at, updated_at, COUNT(*) OVER()
	FROM profile` + formatWhereClause(where) + ` ORDER BY id ASC` + formatLimitOffset(filter.Limit, filter.Offset)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	profiles := make([]*fwt.Profile, 0)
	for rows.Next() {
		var profile fwt.Profile
		if err := rows.Scan(
			&profile.ID,
			&profile.UserID,
			&profile.FirstName,
			&profile.LastName,
			&profile.DateOfBirth,
			&profile.Gender,
			&profile.Height,
			&profile.Weight,
			(*NullTime)(&profile.CreatedAt),
			(*NullTime)(&profile.UpdatedAt),
			&n,
		); err != nil {
			return nil, n, err
		}

		profiles = append(profiles, &profile)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return profiles, n, nil
}
