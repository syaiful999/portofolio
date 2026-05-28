package outsource

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type TransactOutsource struct {
	ID             uuid.UUID      `json:"id"`
	OutsourceId    string         `json:"outsource_id"`
	OutsourceName  string         `json:"outsource_name"`
	DepartmentId   string         `json:"department_id"`
	DepartmentName string         `json:"department_name"`
	Email          sql.NullString `json:"email"`
	Contact        sql.NullString `json:"contact"`
	IsActive       bool           `json:"is_active"`
	CreatedBy      string         `json:"created_by"`
	CreatedDate    time.Time      `json:"created_date"`
	ModifiedBy     string         `json:"modified_by"`
	ModifiedDate   time.Time      `json:"modified_date"`
}

type MasterOutsource struct {
	ID            uuid.UUID `json:"id"`
	OutsourceName string    `json:"outsource_name"`
	Email         string    `json:"email"`
	Contact       string    `json:"contact"`
	IsActive      bool      `json:"is_active"`
	CreatedBy     string    `json:"created_by"`
	CreatedDate   time.Time `json:"created_date"`
	ModifiedBy    string    `json:"modified_by"`
	ModifiedDate  time.Time `json:"modified_date"`
}

type IOutsourceRepository interface {
	GetOutsource(ctx context.Context, skip, take int32, filter, sort string) ([]TransactOutsource, int32, error)
	GetOutsourceDistinctName(ctx context.Context, skip, take int32, filter, sort string) ([]TransactOutsource, int32, error)
	GetOutsourceList() ([]TransactOutsource, error)
	GetOutsourceByTransactId(ctx context.Context, id uuid.UUID) (TransactOutsource, error)
	CreateMasterOutsource(ctx context.Context, value MasterOutsource) (MasterOutsource, error)
	CreateTransactOutsource(ctx context.Context, value TransactOutsource) (TransactOutsource, error)
	UpdateOutsource(ctx context.Context, value MasterOutsource) (MasterOutsource, error)
	UpdateTransactOutsource(ctx context.Context, value TransactOutsource) (TransactOutsource, error)
	DeleteTransactOutsource(ctx context.Context, id, modifiedBy uuid.UUID) error
	GetIDOutsourceByName(ctx context.Context, outsourceName string) (string, error)
	GetIDTransactOutsourceSynced(ctx context.Context, idOutsource, idDepartment string) (string, error)
}

type repository struct{ db *sql.DB }

func NewOutsourceRepository(db *sql.DB) IOutsourceRepository { return &repository{db: db} }

func (r *repository) GetOutsource(ctx context.Context, skip, take int32, filter, sort string) ([]TransactOutsource, int32, error) {
	return nil, 0, nil
}
func (r *repository) GetOutsourceDistinctName(ctx context.Context, skip, take int32, filter, sort string) ([]TransactOutsource, int32, error) {
	return nil, 0, nil
}
func (r *repository) GetOutsourceList() ([]TransactOutsource, error) { return nil, nil }
func (r *repository) GetOutsourceByTransactId(ctx context.Context, id uuid.UUID) (TransactOutsource, error) {
	return TransactOutsource{}, nil
}
func (r *repository) CreateMasterOutsource(ctx context.Context, value MasterOutsource) (MasterOutsource, error) {
	return MasterOutsource{}, nil
}
func (r *repository) CreateTransactOutsource(ctx context.Context, value TransactOutsource) (TransactOutsource, error) {
	return TransactOutsource{}, nil
}
func (r *repository) UpdateOutsource(ctx context.Context, value MasterOutsource) (MasterOutsource, error) {
	return MasterOutsource{}, nil
}
func (r *repository) UpdateTransactOutsource(ctx context.Context, value TransactOutsource) (TransactOutsource, error) {
	return TransactOutsource{}, nil
}
func (r *repository) DeleteTransactOutsource(ctx context.Context, id, modifiedBy uuid.UUID) error {
	return nil
}
func (r *repository) GetIDOutsourceByName(ctx context.Context, outsourceName string) (string, error) {
	return "", nil
}
func (r *repository) GetIDTransactOutsourceSynced(ctx context.Context, idOutsource, idDepartment string) (string, error) {
	return "", nil
}
