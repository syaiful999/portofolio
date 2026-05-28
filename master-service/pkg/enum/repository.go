package enum

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type IEnumRepository interface {
	GetEnum(ctx context.Context, skip, take int32, filter, sort string) ([]MasterEnum, int32, error)
	GetEnumById(ctx context.Context, id uuid.UUID) (MasterEnum, error)
	// returning master_enum where enum type in ('jabatan', 'grade', 'status_contract', 'cost_center', 'gl_number')
	GetEnumExcel() ([]MasterEnum, error)
	UpdateEnum(ctx context.Context, value MasterEnum) (MasterEnum, error)
}

type repository struct {
	db *sql.DB
}

func NewEnumRepository(db *sql.DB) *repository {

	return &repository{
		db: db,
	}
}

func (r *repository) GetEnum(ctx context.Context, skip, take int32, filter, sort string) ([]MasterEnum, int32, error) {

	var items []MasterEnum
	var countData int32

	query := fmt.Sprintf(`SELECT id, is_active, created_by, created_date, modified_by, modified_date, enum_value, enum_type, enum_code
	FROM master.master_enum %s %s limit $1 offset $2;`,
		filter, sort)

	rows, err := r.db.Query(query, take, skip)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var result MasterEnum
		if err := rows.Scan(&result.ID, &result.IsActive, &result.CreatedBy, &result.CreatedDate, &result.ModifiedBy, &result.ModifiedDate, &result.EnumValue, &result.EnumType, &result.EnumCode); err != nil {
			return nil, 0, err
		}
		items = append(items, result)
	}
	if err := rows.Close(); err != nil {
		return nil, 0, err
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	//count data
	queryCount := fmt.Sprintf(`select count(id) from master.master_enum %s;`, filter)

	row := r.db.QueryRow(queryCount)
	errCount := row.Scan(&countData)
	if errCount != nil {
		return nil, 0, err
	}

	return items, countData, nil
}

func (r *repository) GetEnumById(ctx context.Context, id uuid.UUID) (MasterEnum, error) {

	var result MasterEnum

	query := `SELECT id, is_active, created_by, created_date, modified_by, modified_date, enum_value, enum_type
				FROM master.master_enum
				WHERE id=$1`

	row := r.db.QueryRow(query, id)
	err := row.Scan(
		&result.ID,
		&result.IsActive,
		&result.CreatedBy,
		&result.CreatedDate,
		&result.ModifiedBy,
		&result.ModifiedDate,
		&result.EnumValue,
		&result.EnumType,
	)
	return result, err
}

func (r *repository) UpdateEnum(ctx context.Context, value MasterEnum) (MasterEnum, error) {

	var result MasterEnum

	query := `UPDATE master.master_enum SET 
				enum_value=$1,
				modified_by=$2, 
				modified_date=now() 
				WHERE enum_type=$3
			  RETURNING id`
	row := r.db.QueryRow(query,
		value.EnumValue,
		value.ModifiedBy,
		value.EnumType,
	)
	err := row.Scan(
		&result.ID,
	)
	return result, err
}

func (r *repository) GetEnumExcel() ([]MasterEnum, error) {

	var items []MasterEnum

	rows, err := r.db.Query(`
	SELECT id, enum_type , enum_value  
	FROM master.master_enum me 
	WHERE enum_type IN
	('grade', 'status_ptkp', 'status_contract', 'cost_center', 'gl_number', 'gender',
	 'status_marriage', 'blood_type', 'status_contract', 'religion') 
		AND is_active 
	ORDER BY 2,3;`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var result MasterEnum
		if err := rows.Scan(&result.ID, &result.EnumType, &result.EnumValue); err != nil {
			return nil, err
		}
		items = append(items, result)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
