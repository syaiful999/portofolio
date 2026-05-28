package menu

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type IMenuRepository interface {
	GetMenu(ctx context.Context, skip, take int32, filter, sort string) ([]MasterMenu, int32, error)
	GetMenuById(ctx context.Context, id uuid.UUID) ([]MasterMenu, error)
}

type repository struct {
	db *sql.DB
}

func NewMenuRepository(db *sql.DB) *repository {

	return &repository{
		db: db,
	}
}

func (r *repository) GetMenu(ctx context.Context, skip, take int32, filter, sort string) ([]MasterMenu, int32, error) {

	var items []MasterMenu
	var countData int32

	query := fmt.Sprintf(`SELECT id, menu_code, menu_description, menu_path, is_active, created_by, created_date, modified_by, modified_date, id_parent, level
	FROM master.v_master_menu %s %s limit $1 offset $2;`,
		filter, sort)

	rows, err := r.db.Query(query, take, skip)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var result MasterMenu
		if err := rows.Scan(
			&result.ID,
			&result.MenuCode,
			&result.MenuDescription,
			&result.MenuPath,
			&result.IsActive,
			&result.CreatedBy,
			&result.CreatedDate,
			&result.ModifiedBy,
			&result.ModifiedDate,
			&result.IdParent,
			&result.Level,
		); err != nil {
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
	queryCount := fmt.Sprintf(`select count(id) from master.v_master_menu %s;`,
		filter,
	)

	row := r.db.QueryRow(queryCount)
	errCount := row.Scan(&countData)
	if errCount != nil {
		return nil, 0, err
	}

	return items, countData, nil
}

func (r *repository) GetMenuById(ctx context.Context, id uuid.UUID) ([]MasterMenu, error) {

	var items []MasterMenu

	query := `SELECT id, menu_code, menu_description, menu_path, is_active, created_by, created_date, modified_by, modified_date, id_parent, level
			  FROM master.v_master_menu
			  WHERE id=$1 OR id_parent=$1;`

	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var result MasterMenu
		if err := rows.Scan(
			&result.ID,
			&result.MenuCode,
			&result.MenuDescription,
			&result.MenuPath,
			&result.IsActive,
			&result.CreatedBy,
			&result.CreatedDate,
			&result.ModifiedBy,
			&result.ModifiedDate,
			&result.IdParent,
			&result.Level,
		); err != nil {
			return nil, err
		}
		items = append(items, result)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
