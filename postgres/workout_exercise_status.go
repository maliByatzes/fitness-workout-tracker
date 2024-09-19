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
	return nil
}

func (s *WEStatusService) UpdateWEStatus(ctx context.Context, id uint, upd fwt.WEStatusUpdate) (*fwt.WEStatus, error) {
	return nil, nil
}

func (s *WEStatusService) DeleteWEStatus(ctx context.Context, id uint) error {
	return nil
}
