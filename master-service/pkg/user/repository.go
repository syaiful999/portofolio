package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type IUserRepository interface {

	// returning data master_user with specific needs
	// | skip is offset & take is limit in postgres
	// | filter is specific query you want to read
	// | sort is order by
	//
	// GetUser(ctx, 0, 10, "WHERE is_active=TRUE", "ORDER BY created_date")
	GetUser(ctx context.Context, skip, take int32, filter, sort string) ([]MasterUser, int32, error)

	// returning data master_user with specific id (master_user.id)
	//
	// GetUserById(ctx, 1)
	GetUserById(ctx context.Context, id uuid.UUID) (MasterUser, error)
	GetUserByEmail(ctx context.Context, email string) (MasterUser, string, error)

	// create new data master_user
	//
	// newMasterUser MasterUser{ ... }
	//
	// CreateUser(ctx, newMasterUser)
	CreateUser(ctx context.Context, value MasterUser) (MasterUser, string, error)

	// update data master_user by id
	//
	// currentMasterUser MasterUser{ id, ... }
	//
	// UpdateUser(ctx, currentMasterUser)
	UpdateUser(ctx context.Context, value MasterUser) (MasterUser, string, error)
	// update data master_user.password by id
	//
	// currentMasterUser MasterUser{ id, ... }
	//
	// UpdateUserPassword(ctx, currentMasterUser)
	UpdateUserPassword(ctx context.Context, value MasterUser, redirectId string) (MasterUser, error)
	// delete data master_user by id
	//
	// currentMasterUser MasterUser{ id, ... }
	//
	// DeleteUser(ctx, uuid(123), uuid(123))
	DeleteUser(ctx context.Context, id, modifiedBy uuid.UUID) (err error)

	// counting to check if email is already exist or not. only not-deleted user will be checked
	//
	// GetCountUniqueUser(ctx, dataUser)
	GetCountUniqueUser(ctx context.Context, data MasterUser) (int, error)
	GetUserGroupbyRole(ctx context.Context, skip, take int32, filter, sort string) ([]MasterUserGroupbyRole, error)

	// activate data master_user by id from token in header
	//
	// currentMasterUser MasterUser{ id, ... }
	//
	// ActivateUser(ctx, currentMasterUser)
	ActivateUser(ctx context.Context, value MasterUser) (MasterUser, error)

	GetDepartementByOutsourceId(id string) ([]string, error)

	GetUserPasswordHash(ctx context.Context, userId uuid.UUID) (string, error)
}

type repository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *repository {

	return &repository{
		db: db,
	}
}

func (r *repository) GetUser(ctx context.Context, skip, take int32, filter, sort string) ([]MasterUser, int32, error) {

	var items []MasterUser
	var countData int32

	query := fmt.Sprintf(`SELECT id, is_active, created_by, created_date, modified_by, modified_date,
	user_name, name, email, role_id, role_description, role_code, department_id, department_name, 
	outsource_id, outsource_name, location_id, location_name, picture, is_status
	FROM master.v_master_user %s %s limit $1 offset $2;`,
		filter, sort)

	rows, err := r.db.Query(query, take, skip)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var result MasterUser
		if err := rows.Scan(
			&result.ID, &result.IsActive, &result.CreatedBy, &result.CreatedDate, &result.ModifiedBy, &result.ModifiedDate,
			&result.UserName, &result.Name, &result.Email, &result.RoleId, &result.RoleDescription, &result.RoleCode, &result.DepartmentId, &result.DepartmentName,
			&result.OutsourceId, &result.OutsourceName, &result.LocationId, &result.LocationName, &result.Picture, &result.IsStatus,
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
	queryCount := fmt.Sprintf(`select count(id) from master.v_master_user %s;`,
		filter,
	)

	row := r.db.QueryRow(queryCount)
	errCount := row.Scan(&countData)
	if errCount != nil {
		return nil, 0, err
	}

	return items, countData, nil
}

func (r *repository) GetUserById(ctx context.Context, id uuid.UUID) (MasterUser, error) {

	var result MasterUser

	query := `SELECT id, is_active, created_by, created_date, modified_by, modified_date,
				user_name, name, email, role_id, role_description, role_code, department_id, department_name,
				outsource_id, outsource_name, location_id, location_name, picture, is_status
			FROM master.v_master_user
			WHERE id=$1`

	row := r.db.QueryRow(query, id)
	err := row.Scan(
		&result.ID, &result.IsActive, &result.CreatedBy, &result.CreatedDate, &result.ModifiedBy, &result.ModifiedDate,
		&result.UserName, &result.Name, &result.Email, &result.RoleId, &result.RoleDescription, &result.RoleCode, &result.DepartmentId, &result.DepartmentName,
		&result.OutsourceId, &result.OutsourceName, &result.LocationId, &result.LocationName, &result.Picture, &result.IsStatus,
	)
	return result, err
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (MasterUser, string, error) {

	var result MasterUser
	tx, err := r.db.Begin()

	if err != nil {
		tx.Rollback()
		return result, "", err
	}
	query := `SELECT id, user_name, is_status
			FROM master.v_master_user
			WHERE email=$1`

	row := tx.QueryRow(query, email)
	err = row.Scan(&result.ID, &result.UserName, &result.IsStatus)
	if err != nil {
		tx.Rollback()
		return result, "", err
	}

	var idRedirectPage string
	idRedirectPage, err = createRedirectPage(tx, result.ID.String(), "forgot-password")
	if err != nil {
		tx.Rollback()
		return result, "", err
	}
	tx.Commit()

	return result, idRedirectPage, err
}
func (r *repository) CreateUser(ctx context.Context, value MasterUser) (MasterUser, string, error) {

	var result MasterUser

	tx, err := r.db.Begin()
	if err != nil {
		tx.Rollback()
		return result, "", err
	}
	mutation := `INSERT INTO master.master_user
			(created_by, modified_by, user_name, name,  email, role_id, department_id, outsource_id, location_id, picture)
			VALUES($1,$2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING id,  created_date, modified_date, user_name`

	row := tx.QueryRow(mutation,
		value.CreatedBy, value.ModifiedBy,
		value.UserName,
		value.Name,
		value.Email,
		value.RoleId,
		value.DepartmentId,
		value.OutsourceId,
		value.LocationId,
		value.Picture,
	)

	err = row.Scan(
		&result.ID,
		&result.CreatedDate,
		&result.ModifiedDate,
		&result.UserName,
	)

	if err != nil {
		tx.Rollback()
		return result, "", err
	}
	var idRedirectPage string
	idRedirectPage, err = createRedirectPage(tx, result.ID.String(), "generate-password")

	if err != nil {
		tx.Rollback()
		return result, "", err
	}

	tx.Commit()
	return result, idRedirectPage, err

}

func (r *repository) UpdateUser(ctx context.Context, value MasterUser) (MasterUser, string, error) {

	var result MasterUser

	tx, err := r.db.Begin()
	if err != nil {
		tx.Rollback()
		return result, "", err
	}
	query := `UPDATE master.master_user
				SET
					name = $1,
					modified_by = $2,
					modified_date =now() ,
					user_name = $3,
					email=$4,
					role_id=$5,
					department_id=$6,
					outsource_id=$7,
					location_id=$8,
					picture=$9		
				WHERE
					id = $10
				RETURNING id,  modified_by, modified_date, user_name, is_status`
	row := r.db.QueryRow(query,
		value.Name,
		value.ModifiedBy,
		value.UserName,
		value.Email,
		value.RoleId,
		value.DepartmentId,
		value.OutsourceId,
		value.LocationId,
		value.Picture,
		value.ID,
	)
	err = row.Scan(
		&result.ID,
		&result.ModifiedBy,
		&result.ModifiedDate,
		&result.UserName,
		&result.IsStatus,
	)

	if err != nil {
		tx.Rollback()
		return result, "", err
	}

	var idRedirectPage string
	idRedirectPage, err = createRedirectPage(tx, result.ID.String(), "forgot-password")

	if err != nil {
		tx.Rollback()
		return result, "", err
	}
	tx.Commit()

	return result, idRedirectPage, err
}
func (r *repository) UpdateUserPassword(ctx context.Context, value MasterUser, redirectId string) (MasterUser, error) {
	var result MasterUser

	tx, err := r.db.Begin()
	if err != nil {
		tx.Rollback()
		return result, err
	}
	query := `UPDATE master.master_user
		SET
			password = $1,
			password_last_changed = now(),
			password_expiry_date = now() + (
				SELECT (enum_value || ' days')::interval
				FROM master.master_enum
				WHERE enum_type = 'PASSWORD_POLICY'
				  AND enum_code = 'EXPIRY_DAYS'
				  AND is_active = true
				LIMIT 1
			),
			must_change_password = false,
			failed_login_attempts = 0,
			last_failed_login = NULL,
			account_locked_until = NULL,
			modified_by = $2,
			modified_date = now()
		WHERE id = $2
		RETURNING id, modified_by, modified_date, user_name;`
	row := tx.QueryRow(query,
		value.Password,
		value.ID,
	)
	err = row.Scan(
		&result.ID,
		&result.ModifiedBy,
		&result.ModifiedDate,
		&result.UserName,
	)

	if err != nil {
		tx.Rollback()
		return result, err
	}
	err = disableRedirectPage(tx, value.ID.String(), redirectId)
	if err != nil {
		tx.Rollback()
		return result, err
	}
	tx.Commit()
	return result, err
}

func (r *repository) DeleteUser(ctx context.Context, id, modifiedBy uuid.UUID) error {

	query := `UPDATE master.master_user 
		SET modified_date=now(), 
			modified_by = $1, 
			is_active=FALSE
		WHERE id = $2`
	_, err := r.db.Exec(query, modifiedBy, id)
	return err
}

func (r *repository) GetCountUniqueUser(ctx context.Context, dataUser MasterUser) (int, error) {

	var count int

	query := `SELECT count(*) from master.master_user 
		WHERE (LOWER(email) = LOWER($1) OR LOWER(user_name) = LOWER($2))
			AND id != $3
			AND is_active = TRUE`
	row := r.db.QueryRow(query, dataUser.Email, dataUser.UserName, dataUser.ID)
	err := row.Scan(&count)

	return count, err
}

func (r *repository) GetUserGroupbyRole(ctx context.Context, skip, take int32, filter, sort string) ([]MasterUserGroupbyRole, error) {

	var items []MasterUserGroupbyRole

	query := fmt.Sprintf(`
	SELECT
	 role_id, role_code, role_description, "count"
	FROM master.v_master_user_groupby_role  %s %s limit $1 offset $2;`,
		filter, sort)

	rows, err := r.db.Query(query, take, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var result MasterUserGroupbyRole
		if err := rows.Scan(&result.RoleID, &result.RoleCode, &result.RoleDescription, &result.Count); err != nil {
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

	return items, err
}

func (r *repository) ActivateUser(ctx context.Context, value MasterUser) (MasterUser, error) {

	var result MasterUser = value
	query := `UPDATE master.master_user
				SET 
					is_status=TRUE,
					modified_date=now(),
					modified_by=$1
				WHERE
					id = $1
				RETURNING modified_date`
	row := r.db.QueryRow(query, value.ID)
	err := row.Scan(
		&result.ModifiedDate,
	)
	return result, err
}

func createRedirectPage(tx *sql.Tx, idUser, key string) (string, error) {
	var idRedirectPage string
	mutation := `INSERT INTO transaction.transact_redirect_page 
			(created_by, modified_by, page_name)
		VALUES($1, $1, $2)
		RETURNING id`
	row := tx.QueryRow(mutation, idUser, key)

	err := row.Scan(&idRedirectPage)

	return idRedirectPage, err

}

func disableRedirectPage(tx *sql.Tx, userId, redirectId string) error {

	mutation := `
	UPDATE transaction.transact_redirect_page 
	SET 
		is_active=FALSE, 
		modified_date= now(),
		modified_by= $1
	WHERE id=$2
		RETURNING id`
	_, err := tx.Exec(mutation, userId, redirectId)

	return err
}

func (r *repository) GetDepartementByOutsourceId(id string) ([]string, error) {
	var result []string

	query := `SELECT department_id FROM "transaction".transact_outsource where 1=1 and outsource_id=$1 and is_active = true`
	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item string
		if err := rows.Scan(
			&item,
		); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *repository) GetUserPasswordHash(ctx context.Context, userId uuid.UUID) (string, error) {
	var hashedPassword string

	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	query := `SELECT COALESCE(password, '') AS password FROM master.master_user WHERE id = $1`
	err = tx.QueryRowContext(ctx, query, userId).Scan(&hashedPassword)
	if err != nil {
		return "", err
	}

	tx.Commit()
	return hashedPassword, nil
}
