package role

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type IRoleRepository interface {
	GetRole(ctx context.Context, skip, take int32, filter, sort string) ([]MasterRole, int32, error)
	GetRoleById(ctx context.Context, id uuid.UUID) (MasterRole, error)
	GetRoleMapping(id_role string) ([]AuthorizationModel, int32, error)
	GetRoleMappingEmployee(id_role string) ([]AuthorizationModel, error)
	CreateRole(tx *sql.Tx, value MasterRole) (MasterRole, error)
	CreateTransactionRoleMapping(tx *sql.Tx, value AuthorizationModel) error
	UpdateRole(tx *sql.Tx, value MasterRole) (MasterRole, error)
	DeleteRole(tx *sql.Tx, id, modifiedBy uuid.UUID) error
	GetMasterRoleByName(tx *sql.Tx, role_code string) (int32, error)
	GetMasterRoleValidate(tx *sql.Tx, role_code string, id uuid.UUID) (int32, error)
	DeleteOldTransaction(tx *sql.Tx, id uuid.UUID) error
	GetDB() *sql.DB
}

type repository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *repository {

	return &repository{
		db: db,
	}
}

func (r *repository) GetRole(ctx context.Context, skip, take int32, filter, sort string) ([]MasterRole, int32, error) {

	var items []MasterRole
	var countData int32

	query := fmt.Sprintf(`SELECT id, role_code, role_description, is_active, created_by, created_date, modified_by, modified_date
	FROM master.v_master_role %s %s limit $1 offset $2;`,
		filter, sort)

	rows, err := r.db.Query(query, take, skip)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var result MasterRole
		if err := rows.Scan(
			&result.ID,
			&result.RoleCode,
			&result.RoleDescription,
			&result.IsActive,
			&result.CreatedBy,
			&result.CreatedDate,
			&result.ModifiedBy,
			&result.ModifiedDate,
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
	queryCount := fmt.Sprintf(`select count(id) from master.v_master_role %s;`,
		filter,
	)

	row := r.db.QueryRow(queryCount)
	errCount := row.Scan(&countData)
	if errCount != nil {
		return nil, 0, err
	}

	return items, countData, nil
}

func (r *repository) GetRoleMapping(id_role string) ([]AuthorizationModel, int32, error) {

	var items []AuthorizationModel
	var countData int32

	query := fmt.Sprintf(`SELECT id, id_role, id_menu, menu_code, menu_description, menu_path, is_writable, is_active, created_by, created_date, modified_by, modified_date, id_parent, level
	FROM transaction.v_transact_role_mapping where id_role = '%s';`,
		id_role,
	)

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var result AuthorizationModel
		if err := rows.Scan(
			&result.Id,
			&result.IdRole,
			&result.IdMenu,
			&result.MenuCode,
			&result.MenuDescription,
			&result.MenuPath,
			&result.IsWritable,
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
	queryCount := fmt.Sprintf(`SELECT count(a.id)
			FROM transaction.v_transact_role_mapping a
			WHERE a.id_role = '%s'
			AND NOT EXISTS (
				SELECT 1
				FROM transaction.v_transact_role_mapping b
				WHERE b.id_parent = a.id_menu
			);
		`,
		id_role,
	)

	row := r.db.QueryRow(queryCount)
	errCount := row.Scan(&countData)
	if errCount != nil {
		return nil, 0, err
	}

	return items, countData, nil
}

func (r *repository) GetRoleById(ctx context.Context, id uuid.UUID) (MasterRole, error) {

	var result MasterRole

	query := `SELECT id, role_code, role_description, is_active, created_by, created_date, modified_by, modified_date
			  FROM master.v_master_role
			  WHERE id=$1`

	row := r.db.QueryRow(query, id)
	err := row.Scan(
		&result.ID,
		&result.RoleCode,
		&result.RoleDescription,
		&result.IsActive,
		&result.CreatedBy,
		&result.CreatedDate,
		&result.ModifiedBy,
		&result.ModifiedDate,
	)
	return result, err
}

func (r *repository) CreateRole(tx *sql.Tx, value MasterRole) (MasterRole, error) {

	var result MasterRole

	query := `INSERT INTO master.master_role
				(id, role_code, role_description, is_active, created_by, created_date, modified_by, modified_date)
			  VALUES($1, $2, $3, true, $4, now(), $5, now())
			  RETURNING id`

	row := tx.QueryRow(query,
		value.ID,
		value.RoleCode,
		value.RoleDescription,
		value.CreatedBy,
		value.ModifiedBy,
	)

	err := row.Scan(
		&result.ID,
	)
	return result, err

}

func (r *repository) CreateTransactionRoleMapping(tx *sql.Tx, value AuthorizationModel) error {

	query := `INSERT INTO transaction.transact_role_mapping
				(id, id_role, id_menu, is_writable, is_active, created_by, created_date, modified_by, modified_date)
			  VALUES($1, $2, $3, $4, true, $5, now(), $6, now())
			  RETURNING id`

	row := r.db.QueryRow(query,
		value.Id,
		value.IdRole,
		value.IdMenu,
		value.IsWritable,
		value.CreatedBy,
		value.ModifiedBy,
	)

	err := row.Scan(&value.Id)
	return err

}

func (r *repository) UpdateRole(tx *sql.Tx, value MasterRole) (MasterRole, error) {

	var result MasterRole

	query := `UPDATE master.master_role SET 
					role_code=$1, 
					role_description=$2,
					modified_by=$3, 
					modified_date=now()
				WHERE id = $4
				RETURNING id`
	row := r.db.QueryRow(query,
		value.RoleCode,
		value.RoleDescription,
		value.ModifiedBy,
		value.ID,
	)
	err := row.Scan(
		&result.ID,
	)
	return result, err
}

func (r *repository) DeleteRole(tx *sql.Tx, id, modifiedBy uuid.UUID) error {

	query := `UPDATE master.master_role SET modified_date=now(), modified_by = $1, is_active=false WHERE id = $2`
	_, err := r.db.Exec(query, modifiedBy, id)
	return err
}

func (r *repository) GetMasterRoleByName(tx *sql.Tx, role_code string) (int32, error) {

	var count int32

	query := `SELECT count(id) from master.master_role WHERE lower(role_code) = lower($1) and is_active = true`
	row := tx.QueryRow(query, role_code)
	err := row.Scan(&count)

	return count, err
}

func (r *repository) GetMasterRoleValidate(tx *sql.Tx, role_code string, id uuid.UUID) (int32, error) {

	var count int32

	query := `SELECT count(id) from master.master_role WHERE lower(role_code) = lower($1) and is_active = true and id != $2`
	row := r.db.QueryRow(query, role_code, id)
	err := row.Scan(&count)

	return count, err
}

func (r *repository) DeleteOldTransaction(tx *sql.Tx, role_id uuid.UUID) error {
	query := `DELETE FROM transaction.transact_role_mapping WHERE id_role=$1;`
	_, err := tx.Exec(query, role_id)
	return err
}

func (r *repository) GetDB() *sql.DB {
	return r.db
}

func (r *repository) GetRoleMappingEmployee(id_role string) ([]AuthorizationModel, error) {
	var items []AuthorizationModel
	rows, err := r.db.Query(`SELECT 
		 id, id_role, id_menu, menu_code, menu_description, menu_path, is_writable, id_parent, menu_code_parent
	FROM transaction.v_transact_role_mapping 
	WHERE menu_code in ('employment','job-description','address','profile') AND
	menu_code_parent ='Employee' and level=3 and id_role = $1;`, id_role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var result AuthorizationModel
		if err := rows.Scan(
			&result.Id,
			&result.IdRole,
			&result.IdMenu,
			&result.MenuCode,
			&result.MenuDescription,
			&result.MenuPath,
			&result.IsWritable,
			&result.IdParent,
			&result.MenuCodeParent,
		); err != nil {
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
