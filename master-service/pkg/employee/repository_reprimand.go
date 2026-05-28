package employee

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func (r *repository) GetEmployeeReprimand(ctx context.Context, employeeId uuid.UUID) ([]ReprimandModel, error) {

	var items []ReprimandModel

	query := `SELECT
				id, attachment_url, warning_level_id, warning_level_value, valid_time, document_number, description,
				start_date, end_date
			FROM transaction.v_transact_employee_reprimand
			WHERE is_active = TRUE AND employee_id = $1
			 ORDER BY id`

	rows, err := r.db.Query(query, employeeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var result ReprimandModel
		if err := rows.Scan(
			&result.Id, &result.AttachmentUrl, &result.WarningLevelId, &result.WarningLevelValue, &result.ValidTime, &result.DocumentNumber, &result.Description,
			&result.StartDate, &result.EndDate,
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

func (r *repository) InsertEmployeeReprimand(tx *sql.Tx, userId, employeeId string, value ReprimandModel) error {
	query := `
		INSERT INTO transaction.transact_employee_reprimand (
			modified_by, created_by, employee_id, start_date, end_date,
			attachment_url, warning_level_id, valid_time, document_number, description)
		VALUES ($1, $1, $2, $3, $4, 
		$5, $6, $7, $8, $9);`

	_, err := tx.Exec(query, userId, employeeId, value.StartDate, value.EndDate,
		value.AttachmentUrl, value.WarningLevelId, value.ValidTime, value.DocumentNumber, value.Description,
	)

	if err != nil {
		return fmt.Errorf("failed to insert employee reprimand: %v", err)
	}

	return nil
}

func (r *repository) UpdateEmployeeReprimand(tx *sql.Tx, value ReprimandModel) error {
	query := `UPDATE transaction.transact_employee_reprimand 
		SET modified_date=now(), 
		 modified_by =  $1,
			attachment_url =  $2,
			warning_level_id =  $3,
			valid_time =  $4,
			document_number =  $5,
			description = $6,
			start_date = $7,
			end_date = $8
		WHERE id = $9`
	_, err := tx.Exec(query,
		value.ModifiedBy,
		value.AttachmentUrl, value.WarningLevelId, value.ValidTime,
		value.DocumentNumber,
		value.Description,
		value.StartDate,
		value.EndDate,
		value.Id)

	if err != nil {
		return fmt.Errorf("failed to UpdateEmployeeReprimand employee dependents: %v", err)
	}

	return nil
}

func (r *repository) DeleteEmployeeReprimand(tx *sql.Tx, value ReprimandModel) error {
	query := `UPDATE transaction.transact_employee_reprimand 
		SET modified_date=now(), 
			modified_by = $1, 
			is_active=FALSE 
		WHERE id = $2`
	_, err := tx.Exec(query, value.ModifiedBy, value.Id)
	if err != nil {
		return fmt.Errorf("failed to DeleteEmployeeReprimand employee dependents: %v", err)
	}

	return nil
}

func (r *repository) DeleteEmployeeReprimandByEmployeeId(tx *sql.Tx, value ReprimandModel) error {
	query := `UPDATE transaction.transact_employee_reprimand 
		SET modified_date=now(), 
			modified_by = $1, 
			is_active=FALSE 
		WHERE employee_id = $2`
	_, err := tx.Exec(query, value.ModifiedBy, value.EmployeeId)
	if err != nil {
		return fmt.Errorf("failed to DeleteEmployeeReprimandByEmployeeId employee Reprimand: %v", err)
	}

	return nil
}
