package jobcompetency

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type MasterJobCompetency struct {
	ID                     uuid.UUID `json:"ID"`
	IsActive               bool      `json:"IsActive"`
	CreatedBy              string    `json:"CreatedBy"`
	CreatedDate            time.Time `json:"CreatedDate"`
	ModifiedBy             string    `json:"ModifiedBy"`
	ModifiedDate           time.Time `json:"ModifiedDate"`
	CompetencyName         string    `json:"CompetencyName"`
	CompetencyCategoryName string    `json:"CompetencyCategoryName"`
	CompetencyCategoryId   string    `json:"CompetencyCategoryId"`
	CompetencyGroup        string    `json:"CompetencyGroup"`
	Description            string    `json:"Description"`
}

type TransactJobCompetency struct {
	ID                     uuid.UUID      `json:"ID"`
	IsActive               bool           `json:"IsActive"`
	CreatedBy              string         `json:"CreatedBy"`
	CreatedDate            time.Time      `json:"CreatedDate"`
	ModifiedBy             string         `json:"ModifiedBy"`
	ModifiedDate           time.Time      `json:"ModifiedDate"`
	CompetencyName         sql.NullString `json:"CompetencyName"`
	CompetencyCategoryName string         `json:"CompetencyCategoryName"`
	JobCompetencyID        sql.NullString `json:"JobCompetencyID"`
	IsRequired             bool           `json:"IsRequired"`
	TotalTraining          sql.NullInt32  `json:"TotalTraining"`
	CompetencyGroup        string         `json:"CompetencyGroup"`
}

type IJobCompetencyRepository interface {
	GetJobCompetency(ctx context.Context, skip, take int32, filter, sort string) ([]MasterJobCompetency, int32, error)
	GetJobCompetencyById(ctx context.Context, id uuid.UUID) (MasterJobCompetency, error)
	GetJobCompetencyGroupByDescription(ctx context.Context, jobDescriptionId uuid.UUID) ([]TransactJobCompetency, error)
	GetJobCompetencyGroupList(ctx context.Context) ([]MasterJobCompetency, error)
	CreateJobCompetency(ctx context.Context, value MasterJobCompetency) (MasterJobCompetency, error)
	UpdateJobCompetency(ctx context.Context, value MasterJobCompetency) (MasterJobCompetency, error)
	DeleteJobCompetency(ctx context.Context, id, modifiedBy uuid.UUID) (err error)
	GetCountJobCompetencyIsSynced(ctx context.Context, jobCompetencyId uuid.UUID) (int, error)
	GetCountJobCompetencyDuplication(ctx context.Context, value MasterJobCompetency) (int, error)
	GetJobCompetencyTrainingEmployee(ctx context.Context, employeeId string) (map[string]int32, error)
	GetJobCompetencyGroup() ([]string, error)
}

type repository struct{ db *sql.DB }

func NewJobCompetencyRepository(db *sql.DB) IJobCompetencyRepository { return &repository{db: db} }

func (r *repository) GetJobCompetency(ctx context.Context, skip, take int32, filter, sort string) ([]MasterJobCompetency, int32, error) {
	return nil, 0, nil
}
func (r *repository) GetJobCompetencyById(ctx context.Context, id uuid.UUID) (MasterJobCompetency, error) {
	return MasterJobCompetency{}, nil
}
func (r *repository) GetJobCompetencyGroupByDescription(ctx context.Context, jobDescriptionId uuid.UUID) ([]TransactJobCompetency, error) {
	return nil, nil
}
func (r *repository) GetJobCompetencyGroupList(ctx context.Context) ([]MasterJobCompetency, error) {
	return nil, nil
}
func (r *repository) CreateJobCompetency(ctx context.Context, value MasterJobCompetency) (MasterJobCompetency, error) {
	return MasterJobCompetency{}, nil
}
func (r *repository) UpdateJobCompetency(ctx context.Context, value MasterJobCompetency) (MasterJobCompetency, error) {
	return MasterJobCompetency{}, nil
}
func (r *repository) DeleteJobCompetency(ctx context.Context, id, modifiedBy uuid.UUID) error {
	return nil
}
func (r *repository) GetCountJobCompetencyIsSynced(ctx context.Context, jobCompetencyId uuid.UUID) (int, error) {
	return 0, nil
}
func (r *repository) GetCountJobCompetencyDuplication(ctx context.Context, value MasterJobCompetency) (int, error) {
	return 0, nil
}
func (r *repository) GetJobCompetencyTrainingEmployee(ctx context.Context, employeeId string) (map[string]int32, error) {
	return nil, nil
}
func (r *repository) GetJobCompetencyGroup() ([]string, error) { return nil, nil }
