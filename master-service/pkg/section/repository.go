package section

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type MasterSection struct {
	ID                   uuid.UUID `json:"ID"`
	IsActive             bool      `json:"IsActive"`
	CreatedBy            string    `json:"CreatedBy"`
	CreatedDate          time.Time `json:"CreatedDate"`
	ModifiedBy           string    `json:"ModifiedBy"`
	ModifiedDate         time.Time `json:"ModifiedDate"`
	SectionName          string    `json:"Name"`
	DepartmentDetailID   uuid.UUID `json:"DepartmentDetailID"`
	DepartmentDetailName string    `json:"DepartmentDetailName"`
	DepartmentID         uuid.UUID `json:"DepartmentID"`
	DepartmentName       string    `json:"DepartmentName"`
}

type ISectionRepository interface {
	GetSection(ctx context.Context, skip, take int32, filter, sort string) ([]MasterSection, int32, error)
	GetSectionList() ([]MasterSection, error)
	GetSectionById(ctx context.Context, id uuid.UUID, departmentId uuid.UUID) (MasterSection, error)
	CreateSection(ctx context.Context, value MasterSection) (MasterSection, error)
	UpdateSection(ctx context.Context, value MasterSection) (MasterSection, error)
	DeleteSection(ctx context.Context, id, modifiedBy uuid.UUID) error
	GetCountSectionByName(ctx context.Context, sectionName string, id uuid.UUID) (int, error)
}

type repository struct{ db *sql.DB }

func NewSectionRepository(db *sql.DB) ISectionRepository { return &repository{db: db} }

func (r *repository) GetSection(ctx context.Context, skip, take int32, filter, sort string) ([]MasterSection, int32, error) {
	return nil, 0, nil
}
func (r *repository) GetSectionList() ([]MasterSection, error) { return nil, nil }
func (r *repository) GetSectionById(ctx context.Context, id uuid.UUID, departmentId uuid.UUID) (MasterSection, error) {
	return MasterSection{}, nil
}
func (r *repository) CreateSection(ctx context.Context, value MasterSection) (MasterSection, error) {
	return MasterSection{}, nil
}
func (r *repository) UpdateSection(ctx context.Context, value MasterSection) (MasterSection, error) {
	return MasterSection{}, nil
}
func (r *repository) DeleteSection(ctx context.Context, id, modifiedBy uuid.UUID) error { return nil }
func (r *repository) GetCountSectionByName(ctx context.Context, sectionName string, id uuid.UUID) (int, error) {
	return 0, nil
}
