package employee

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func (r *repository) GetEmployeeDependent(ctx context.Context, employeeId uuid.UUID) ([]DependentModel, error) {

	var items []DependentModel

	query := `SELECT
				id, dependent_name, birthdate, phonenumber, relation_id, relation_name
			FROM master.v_master_employee_dependent
			WHERE is_active = TRUE AND employee_id = $1
			 ORDER BY id`

	rows, err := r.db.Query(query, employeeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var result DependentModel
		if err := rows.Scan(
			&result.Id, &result.DependentName, &result.Birthdate, &result.Phonenumber, &result.RelationId, &result.RelationName,
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

func (r *repository) InsertEmployeeMultipleDependents(tx *sql.Tx, userId, employeeId string, values []DependentModel) error {
	if len(values) == 0 {
		return nil
	}

	query := `
		INSERT INTO master.master_employee_dependent (
			modified_by, created_by, employee_id, dependent_name, relation_id, birthdate, birthplace, phonenumber)
		VALUES `

	// Prepare parameters
	var params []interface{}
	for i, value := range values {
		if i > 0 {
			query += ", "
		}
		query +=
			fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				i*8+1, i*8+2, i*8+3, i*8+4, i*8+5, i*8+6, i*8+7, i*8+8)

		params = append(params,
			userId,
			userId,
			employeeId,
			value.DependentName,
			value.RelationId,
			value.Birthdate.Time,
			value.Birthplace.String,
			value.Phonenumber.String,
		)
	}
	_, err := tx.Exec(query, params...)
	if err != nil {
		return fmt.Errorf("failed to InsertEmployeeMultipleDependents employee dependents: %v", err)
	}

	return nil
}

func (r *repository) InsertEmployeeDependent(tx *sql.Tx, employeeId string, value DependentModel) error {
	query := `INSERT INTO master.master_employee_dependent 
	(created_by, modified_by, employee_id, dependent_name, relation_id, birthdate, birthplace, phonenumber)
	VALUES($1, $1, $2,
	 $3, $4, $5, $6, $7)`
	_, err := tx.Exec(query,
		value.ModifiedBy,
		employeeId,
		value.DependentName,
		value.RelationId,
		value.Birthdate,
		value.Birthplace,
		value.Phonenumber)

	if err != nil {
		return fmt.Errorf("failed to InsertEmployeeDependent employee dependents: %v", err)
	}

	return nil
}
func (r *repository) UpdateEmployeeDependent(tx *sql.Tx, value DependentModel) error {
	query := `UPDATE master.master_employee_dependent 
		SET modified_date=now(), 
			modified_by = $1, 
			dependent_name = $2,
			relation_id = $3,
			birthdate = $4,
			birthplace = $5,
			phonenumber = $6
		WHERE id = $7`
	_, err := tx.Exec(query,
		value.ModifiedBy,
		value.DependentName,
		value.RelationId,
		value.Birthdate,
		value.Birthplace,
		value.Phonenumber,
		value.Id)

	if err != nil {
		return fmt.Errorf("failed to UpdateEmployeeDependent employee dependents: %v", err)
	}

	return nil
}

func (r *repository) DeleteEmployeeDependent(tx *sql.Tx, value DependentModel) error {
	query := `UPDATE master.master_employee_dependent 
		SET modified_date=now(), 
			modified_by = $1, 
			is_active=FALSE 
		WHERE id = $2`
	_, err := tx.Exec(query, value.ModifiedBy, value.Id)
	if err != nil {
		return fmt.Errorf("failed to DeleteEmployeeDependent employee dependents: %v", err)
	}

	return nil
}

func (r *repository) DeleteEmployeeDependentByEmployeeId(tx *sql.Tx, value DependentModel) error {
	query := `UPDATE master.master_employee_dependent 
		SET modified_date=now(), 
			modified_by = $1, 
			is_active=FALSE 
		WHERE employee_id = $2`
	_, err := tx.Exec(query, value.ModifiedBy, value.EmployeeId)
	if err != nil {
		return fmt.Errorf("failed to DeleteEmployeeDependentByEmployeeId employee dependents: %v", err)
	}

	return nil
}
