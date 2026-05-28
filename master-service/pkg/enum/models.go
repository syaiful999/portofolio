package enum

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type MasterEnum struct {
	ID           uuid.UUID      `json:"ID"`
	IsActive     bool           `json:"IsActive"`
	CreatedBy    string         `json:"CreatedBy"`
	CreatedDate  time.Time      `json:"CreatedDate"`
	ModifiedBy   string         `json:"ModifiedBy"`
	ModifiedDate time.Time      `json:"ModifiedDate"`
	EnumValue    string         `json:"EnumValue"`
	EnumType     string         `json:"EnumType"`
	EnumCode     sql.NullString `json:"EnumCode"`
}
