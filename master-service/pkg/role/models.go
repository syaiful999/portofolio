package role

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type MasterRole struct {
	ID              uuid.UUID      `json:"id" db:"id"`
	RoleCode        string         `json:"role_name" db:"role_name"`
	RoleDescription sql.NullString `json:"role_description" db:"role_description"`
	IsActive        bool           `json:"is_active" db:"is_active"`
	CreatedBy       string         `json:"created_by" db:"created_by"`
	CreatedDate     time.Time      `json:"created_date" db:"created_date"`
	ModifiedBy      string         `json:"modified_by" db:"modified_by"`
	ModifiedDate    time.Time      `json:"modified_date" db:"modified_dates"`
}

type AuthorizationModel struct {
	Id              uuid.UUID      `json:"id" db:"id"`
	IdRole          uuid.UUID      `json:"id_role" db:"id_role"`
	IdMenu          string         `json:"id_menu" db:"mid_menu"`
	MenuCode        string         `json:"menu_code" db:"menu_code"`
	MenuDescription string         `json:"menu_description" db:"menu_description"`
	MenuPath        string         `json:"menu_path" db:"menu_path"`
	IsWritable      bool           `json:"is_writable" db:"is_writable"`
	IsActive        bool           `json:"is_active" db:"is_active"`
	CreatedBy       string         `json:"created_by" db:"created_by"`
	CreatedDate     time.Time      `json:"created_date" db:"created_date"`
	ModifiedBy      string         `json:"modified_by" db:"modified_by"`
	ModifiedDate    time.Time      `json:"modified_date" db:"modified_date"`
	IdParent        sql.NullString `json:"id_parent" db:"id_parent"`
	Level           int            `json:"level" db:"level"`
	MenuCodeParent  sql.NullString `json:"menu_code_parent" db:"menu_code_parent"`
}
