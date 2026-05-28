package shift

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type MasterShift struct {
	ID           uuid.UUID      `json:"id"`
	IsActive     bool           `json:"is_active"`
	CreatedBy    string         `json:"created_by"`
	CreatedDate  interface{}    `json:"created_date"`
	ModifiedBy   string         `json:"modified_by"`
	ModifiedDate interface{}    `json:"modified_date"`
	ShiftName    string         `json:"shift_name"`
	ShiftStart   string         `json:"shift_start"`
	ShiftEnd     string         `json:"shift_end"`
	BreakStart   sql.NullString `json:"break_start"`
	BreakEnd     sql.NullString `json:"break_end"`
	AutoOvertime int32          `json:"auto_overtime"`
	IsOvertime   bool           `json:"is_overtime"`
	Duration     int16          `json:"duration"`
	ShiftTypeId  string         `json:"shift_type_id"`
	ShiftType    string         `json:"shift_type"`
}

type IShiftRepository interface {
	GetShift(ctx context.Context, skip, take int32, filter, sort string) ([]MasterShift, int32, error)
	GetShiftById(ctx context.Context, id uuid.UUID) (MasterShift, error)
	CreateShift(ctx context.Context, value MasterShift) (MasterShift, error)
	UpdateShift(ctx context.Context, value MasterShift) (MasterShift, error)
	DeleteShift(ctx context.Context, id, modifiedBy uuid.UUID) error
	GetCountShiftByName(ctx context.Context, shiftName string, id uuid.UUID) (int, error)
}

type repository struct{ db *sql.DB }

func NewShiftRepository(db *sql.DB) IShiftRepository { return &repository{db: db} }

func (r *repository) GetShift(ctx context.Context, skip, take int32, filter, sort string) ([]MasterShift, int32, error) {
	return nil, 0, nil
}
func (r *repository) GetShiftById(ctx context.Context, id uuid.UUID) (MasterShift, error) {
	return MasterShift{}, nil
}
func (r *repository) CreateShift(ctx context.Context, value MasterShift) (MasterShift, error) {
	return MasterShift{}, nil
}
func (r *repository) UpdateShift(ctx context.Context, value MasterShift) (MasterShift, error) {
	return MasterShift{}, nil
}
func (r *repository) DeleteShift(ctx context.Context, id, modifiedBy uuid.UUID) error { return nil }
func (r *repository) GetCountShiftByName(ctx context.Context, shiftName string, id uuid.UUID) (int, error) {
	return 0, nil
}
