package department

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type MasterDepartment struct {
	ID             uuid.UUID      `json:"ID"`
	IsActive       bool           `json:"IsActive"`
	CreatedBy      string         `json:"CreatedBy"`
	CreatedDate    time.Time      `json:"CreatedDate"`
	ModifiedBy     string         `json:"ModifiedBy"`
	ModifiedDate   time.Time      `json:"ModifiedDate"`
	DepartmentName sql.NullString `json:"Name"`
}

type MasterDepartmentCount struct {
	CountOutsource        int `json:"count_outsource"`
	CountEmployee         int `json:"count_employee"`
	CountDepartmentDetail int `json:"count_department_detail"`
}

type IDepartmentRepository interface {
	GetDepartment(ctx context.Context, skip, take int32, filter, sort string) ([]MasterDepartment, int32, error)
	GetDepartmentList() ([]MasterDepartment, error)
	GetDepartmentById(ctx context.Context, id uuid.UUID) (MasterDepartment, error)
	CreateDepartment(ctx context.Context, value MasterDepartment) (MasterDepartment, error)
	UpdateDepartment(ctx context.Context, value MasterDepartment) (MasterDepartment, error)
	DeleteDepartment(ctx context.Context, id, modifiedBy uuid.UUID) error
	GetCountDepartmentByName(ctx context.Context, departmentName string, id uuid.UUID) (int, error)
	GetCountDepartmentIsSynced(ctx context.Context, departmentId uuid.UUID) (MasterDepartmentCount, error)
}

type repository struct{ db *sql.DB }

func NewDepartmentRepository(db *sql.DB) IDepartmentRepository { return &repository{db: db} }

func (r *repository) GetDepartment(ctx context.Context, skip, take int32, filter, sort string) ([]MasterDepartment, int32, error) {
	return nil, 0, nil
}
func (r *repository) GetDepartmentList() ([]MasterDepartment, error) { return nil, nil }
func (r *repository) GetDepartmentById(ctx context.Context, id uuid.UUID) (MasterDepartment, error) {
	return MasterDepartment{}, nil
}
func (r *repository) CreateDepartment(ctx context.Context, value MasterDepartment) (MasterDepartment, error) {
	return MasterDepartment{}, nil
}
func (r *repository) UpdateDepartment(ctx context.Context, value MasterDepartment) (MasterDepartment, error) {
	return MasterDepartment{}, nil
}
func (r *repository) DeleteDepartment(ctx context.Context, id, modifiedBy uuid.UUID) error {
	return nil
}
func (r *repository) GetCountDepartmentByName(ctx context.Context, departmentName string, id uuid.UUID) (int, error) {
	return 0, nil
}
func (r *repository) GetCountDepartmentIsSynced(ctx context.Context, departmentId uuid.UUID) (MasterDepartmentCount, error) {
	return MasterDepartmentCount{}, nil
}
