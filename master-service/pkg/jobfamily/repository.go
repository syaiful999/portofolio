package jobfamily

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type MasterJobFamily struct {
	ID            uuid.UUID `json:"ID"`
	IsActive      bool      `json:"IsActive"`
	CreatedBy     string    `json:"CreatedBy"`
	CreatedDate   time.Time `json:"CreatedDate"`
	ModifiedBy    string    `json:"ModifiedBy"`
	ModifiedDate  time.Time `json:"ModifiedDate"`
	JobFamilyName string    `json:"JobFamilyName"`
	SequenceNo    int32     `json:"SequenceNo"`
}

type IJobFamilyRepository interface {
	GetJobFamily(ctx context.Context, skip, take int32, filter, sort string) ([]MasterJobFamily, int32, error)
	GetJobFamilyList() ([]MasterJobFamily, error)
	GetJobFamilyById(ctx context.Context, id uuid.UUID) (MasterJobFamily, error)
	CreateJobFamily(ctx context.Context, value MasterJobFamily) (MasterJobFamily, error)
	UpdateJobFamily(ctx context.Context, value MasterJobFamily) (MasterJobFamily, error)
	DeleteJobFamily(ctx context.Context, id, modifiedBy uuid.UUID) error
	GetCountJobFamilyByName(ctx context.Context, jobFamilyName string, id uuid.UUID) (int, error)
	GetCountJobFamilyIsSynced(ctx context.Context, jobFamilyId uuid.UUID) (int, error)
}

type repository struct{ db *sql.DB }

func NewJobFamilyRepository(db *sql.DB) IJobFamilyRepository { return &repository{db: db} }

func (r *repository) GetJobFamily(ctx context.Context, skip, take int32, filter, sort string) ([]MasterJobFamily, int32, error) {
	return nil, 0, nil
}
func (r *repository) GetJobFamilyList() ([]MasterJobFamily, error) { return nil, nil }
func (r *repository) GetJobFamilyById(ctx context.Context, id uuid.UUID) (MasterJobFamily, error) {
	return MasterJobFamily{}, nil
}
func (r *repository) CreateJobFamily(ctx context.Context, value MasterJobFamily) (MasterJobFamily, error) {
	return MasterJobFamily{}, nil
}
func (r *repository) UpdateJobFamily(ctx context.Context, value MasterJobFamily) (MasterJobFamily, error) {
	return MasterJobFamily{}, nil
}
func (r *repository) DeleteJobFamily(ctx context.Context, id, modifiedBy uuid.UUID) error {
	return nil
}
func (r *repository) GetCountJobFamilyByName(ctx context.Context, jobFamilyName string, id uuid.UUID) (int, error) {
	return 0, nil
}
func (r *repository) GetCountJobFamilyIsSynced(ctx context.Context, jobFamilyId uuid.UUID) (int, error) {
	return 0, nil
}
