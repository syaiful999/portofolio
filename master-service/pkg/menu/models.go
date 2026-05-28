package menu

import (
	"time"

	"github.com/google/uuid"
)

type MasterMenu struct {
	ID              uuid.UUID     `json:"id" db:"id"`
	MenuCode        string        `json:"menu_code" db:"menu_code"`
	MenuDescription string        `json:"menu_description" db:"menu_description"`
	MenuPath        string        `json:"menu_path" db:"menu_path"`
	IsActive        bool          `json:"is_active" db:"is_active"`
	CreatedBy       string        `json:"created_by" db:"created_by"`
	CreatedDate     time.Time     `json:"created_date" db:"created_date"`
	ModifiedBy      string        `json:"modified_by" db:"modified_by"`
	ModifiedDate    time.Time     `json:"modified_date" db:"modified_date"`
	IdParent        uuid.NullUUID `json:"id_parent" db:"id_parent"`
	Level           int32         `json:"level" db:"level"`
}
