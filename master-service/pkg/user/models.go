package user

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type MasterUser struct {
	ID           uuid.UUID `json:"id"`
	IsActive     bool      `json:"is_active"`
	CreatedBy    string    `json:"created_by"`
	CreatedDate  time.Time `json:"created_date"`
	ModifiedBy   string    `json:"modified_by"`
	ModifiedDate time.Time `json:"modified_date"`
	UserName     string    `json:"user_name"`
	Name         string    `json:"name"`

	Password        string         `json:"password"`
	Email           string         `json:"email"`
	RoleId          sql.NullString `json:"role_id"`
	RoleDescription sql.NullString `json:"role_description"`
	RoleCode        sql.NullString `json:"role_code"`
	DepartmentId    sql.NullString `json:"department_id"`
	DepartmentName  sql.NullString `json:"department_name"`
	OutsourceId     sql.NullString `json:"outsource_id"`
	OutsourceName   sql.NullString `json:"outsource_name"`
	LocationId      sql.NullString `json:"location_id"`
	LocationName    sql.NullString `json:"location_name"`
	IsStatus        bool           `json:"is_status"`
	Picture         sql.NullString `json:"picture"`
}

type MasterUserGroupbyRole struct {
	RoleID          uuid.UUID `json:"id"`
	RoleCode        string    `json:"role_code"`
	RoleDescription string    `json:"role_description"`
	Count           int32     `json:"count"`
}

type TransactRedirectPage struct {
	ID           uuid.UUID `json:"id"`
	CreatedDate  time.Time `json:"created_date"`
	CreatedBy    string    `json:"created_by"`
	ModifiedDate time.Time `json:"modified_date"`
	ModifiedBy   string    `json:"modified_by"`
	ExpiredDate  time.Time `json:"expired_date"`
	PageName     string    `json:"page_name"`
}
