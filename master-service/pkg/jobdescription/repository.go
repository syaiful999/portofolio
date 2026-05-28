package jobdescription

import (
	"context"
	"database/sql"
	"moyo-master-service/pkg/jobcompetency"
	"time"

	"github.com/google/uuid"
)

type MasterJobDescription struct {
	ID                 uuid.UUID      `json:"ID"`
	IsActive           bool           `json:"IsActive"`
	CreatedBy          string         `json:"CreatedBy"`
	CreatedDate        time.Time      `json:"CreatedDate"`
	ModifiedBy         string         `json:"ModifiedBy"`
	ModifiedDate       time.Time      `json:"ModifiedDate"`
	JobResponsibility  sql.NullString `json:"JobResponsibility"`
	JobSpecification   sql.NullString `json:"JobSpecification"`
	JobDescription     sql.NullString `json:"JobDescription"`
	CompetencyName     sql.NullString `json:"CompetencyName"`
	JobFamilyName      string         `json:"JobFamilyName"`
	JobFamilyId        string         `json:"JobFamilyId"`
	DepartmentName     string         `json:"DepartmentName"`
	DepartmentId       string         `json:"DepartmentId"`
	TotalOtherTraining int32          `json:"TotalOtherTraining"`
}

type IJobDescriptionRepository interface {
	GetJobDescription(ctx context.Context, skip, take int32, filter, sort string) ([]MasterJobDescription, int32, error)
	GetJobDescriptionList() ([]MasterJobDescription, error)
	GetJobDescriptionCompetenciesName(ctx context.Context, skip, take int32, filter, sort string) ([]MasterJobDescription, int32, error)
	GetJobDescriptionById(ctx context.Context, id uuid.UUID) (MasterJobDescription, error)
	CreateJobDescription(tx *sql.Tx, value MasterJobDescription) (MasterJobDescription, error)
	UpdateJobDescription(tx *sql.Tx, value MasterJobDescription) (MasterJobDescription, error)
	DeleteJobDescription(tx *sql.Tx, id, modifiedBy uuid.UUID) error
	GetCountJobDescriptionByName(ctx context.Context, jobDescriptionName string, id uuid.UUID) (int, error)
	GetCountJobDescriptionIsSynced(ctx context.Context, jobDescriptionId uuid.UUID) (int, error)
	CreateJobMultipleCompetencies(ctx context.Context, tx *sql.Tx, competencies []string, args []interface{}) error
	UpdateTransactJobDescription(tx *sql.Tx, value jobcompetency.TransactJobCompetency) error
	DeleteTransactJobDescription(tx *sql.Tx, id, modifiedBy uuid.UUID) error
	DeleteTransactJobDescriptionByJobDescriptionId(tx *sql.Tx, jobDescriptionId, modifiedBy uuid.UUID) error
	GetDB() *sql.DB
}

type repository struct{ db *sql.DB }

func NewJobDescriptionRepository(db *sql.DB) IJobDescriptionRepository { return &repository{db: db} }

func (r *repository) GetJobDescription(ctx context.Context, skip, take int32, filter, sort string) ([]MasterJobDescription, int32, error) {
	return nil, 0, nil
}
func (r *repository) GetJobDescriptionList() ([]MasterJobDescription, error) { return nil, nil }
func (r *repository) GetJobDescriptionCompetenciesName(ctx context.Context, skip, take int32, filter, sort string) ([]MasterJobDescription, int32, error) {
	return nil, 0, nil
}
func (r *repository) GetJobDescriptionById(ctx context.Context, id uuid.UUID) (MasterJobDescription, error) {
	return MasterJobDescription{}, nil
}
func (r *repository) CreateJobDescription(tx *sql.Tx, value MasterJobDescription) (MasterJobDescription, error) {
	return MasterJobDescription{}, nil
}
func (r *repository) UpdateJobDescription(tx *sql.Tx, value MasterJobDescription) (MasterJobDescription, error) {
	return MasterJobDescription{}, nil
}
func (r *repository) DeleteJobDescription(tx *sql.Tx, id, modifiedBy uuid.UUID) error { return nil }
func (r *repository) GetCountJobDescriptionByName(ctx context.Context, jobDescriptionName string, id uuid.UUID) (int, error) {
	return 0, nil
}
func (r *repository) GetCountJobDescriptionIsSynced(ctx context.Context, jobDescriptionId uuid.UUID) (int, error) {
	return 0, nil
}
func (r *repository) CreateJobMultipleCompetencies(ctx context.Context, tx *sql.Tx, competencies []string, args []interface{}) error {
	return nil
}
func (r *repository) UpdateTransactJobDescription(tx *sql.Tx, value jobcompetency.TransactJobCompetency) error {
	return nil
}
func (r *repository) DeleteTransactJobDescription(tx *sql.Tx, id, modifiedBy uuid.UUID) error {
	return nil
}
func (r *repository) DeleteTransactJobDescriptionByJobDescriptionId(tx *sql.Tx, jobDescriptionId, modifiedBy uuid.UUID) error {
	return nil
}
func (r *repository) GetDB() *sql.DB { return r.db }
