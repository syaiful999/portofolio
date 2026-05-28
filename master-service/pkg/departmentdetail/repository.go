package departmentdetail

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type MasterDepartmentDetail struct {
	ID                   uuid.UUID `json:"ID"`
	IsActive             bool      `json:"IsActive"`
	CreatedDate          time.Time `json:"CreatedDate"`
	CreatedBy            string    `json:"CreatedBy"`
	ModifiedDate         time.Time `json:"ModifiedDate"`
	ModifiedBy           string    `json:"ModifiedBy"`
	DepartmentDetailName string    `json:"DepartmentDetailName"`
	DepartmentID         uuid.UUID `json:"DepartmentID"`
	DepartmentName       string    `json:"Department"`
}

type IDepartmentDetailRepository interface {
	GetDepartmentDetail(ctx context.Context, skip, take int32, filter, sort string) ([]MasterDepartmentDetail, int32, error)
	GetDepartmentDetailList() ([]MasterDepartmentDetail, error)
	GetDepartmentDetailByIdAdmin(ctx context.Context, id uuid.UUID) (MasterDepartmentDetail, error)
	GetDepartmentDetailById(ctx context.Context, id, departmentId uuid.UUID) (MasterDepartmentDetail, error)
	CreateDepartmentDetail(ctx context.Context, value MasterDepartmentDetail) (MasterDepartmentDetail, error)
	UpdateDepartmentDetail(ctx context.Context, value MasterDepartmentDetail) (MasterDepartmentDetail, error)
	DeleteDepartmentDetail(ctx context.Context, id, modifiedBy uuid.UUID) error
	GetCountDepartmentDetailByNameDepartment(ctx context.Context, departmentName, department_id string, id uuid.UUID) (int, error)
}

type repository struct{ db *sql.DB }

func NewDepartmentDetailRepository(db *sql.DB) IDepartmentDetailRepository {
	return &repository{db: db}
}

func (r *repository) GetDepartmentDetail(ctx context.Context, skip, take int32, filter, sort string) ([]MasterDepartmentDetail, int32, error) {
	return nil, 0, nil
}
func (r *repository) GetDepartmentDetailList() ([]MasterDepartmentDetail, error) { return nil, nil }
func (r *repository) GetDepartmentDetailByIdAdmin(ctx context.Context, id uuid.UUID) (MasterDepartmentDetail, error) {
	return MasterDepartmentDetail{}, nil
}
func (r *repository) GetDepartmentDetailById(ctx context.Context, id, departmentId uuid.UUID) (MasterDepartmentDetail, error) {
	return MasterDepartmentDetail{}, nil
}
func (r *repository) CreateDepartmentDetail(ctx context.Context, value MasterDepartmentDetail) (MasterDepartmentDetail, error) {
	return MasterDepartmentDetail{}, nil
}
func (r *repository) UpdateDepartmentDetail(ctx context.Context, value MasterDepartmentDetail) (MasterDepartmentDetail, error) {
	return MasterDepartmentDetail{}, nil
}
func (r *repository) DeleteDepartmentDetail(ctx context.Context, id, modifiedBy uuid.UUID) error {
	return nil
}
func (r *repository) GetCountDepartmentDetailByNameDepartment(ctx context.Context, departmentName, department_id string, id uuid.UUID) (int, error) {
	return 0, nil
}
