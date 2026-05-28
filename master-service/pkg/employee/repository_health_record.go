package employee

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func (r *repository) GetEmployeeHealthRecord(ctx context.Context, employeeId uuid.UUID) ([]HealthRecordModel, error) {

	var items []HealthRecordModel

	query := `SELECT 
				id, mcu_date, status_id, status, health_description, mcu_url,
				mcu_followup_url, status_followup_id , status_followup, periode
			FROM transaction.v_transact_employee_health_record
			WHERE is_active = TRUE AND employee_id = $1
			 ORDER BY id`

	rows, err := r.db.Query(query, employeeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var result HealthRecordModel
		if err := rows.Scan(
			&result.Id, &result.McuDate, &result.StatusId, &result.Status, &result.HealthDescription, &result.McuUrl,
			&result.McuFollowupUrl, &result.StatusFollowupId, &result.StatusFollowup, &result.Periode,
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

func (r *repository) InsertEmployeeMultipleHealthRecords(tx *sql.Tx, ctx context.Context, userId, employeeId string, values []HealthRecordModel) error {
	if len(values) == 0 {
		return nil
	}
	query := `
		INSERT INTO transaction.transact_employee_health_record (
			modified_by, created_by, employee_id, 
			mcu_date, periode, status_id, description,
			mcu_url, mcu_followup_url, status_followup_id
		) VALUES `

	var params []interface{}
	for i, value := range values {
		if i > 0 {
			query += ", "
		}
		query +=
			fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				i*10+1, i*10+2, i*10+3, i*10+4, i*10+5, i*10+6, i*10+7, i*10+8, i*10+9, i*10+10)

		params = append(params,
			userId,
			userId,
			employeeId,
			value.McuDate,
			value.Periode,
			value.StatusId,
			value.HealthDescription,
			value.McuUrl,
			value.McuFollowupUrl,
			value.StatusFollowupId,
		)
	}
	_, err := tx.ExecContext(ctx, query, params...)
	if err != nil {
		return fmt.Errorf("failed to insert employee health records: %v", err)
	}

	return nil
}

func (r *repository) InsertEmployeeHealthRecord(tx *sql.Tx, userId, employeeId string, value HealthRecordModel) error {
	query := `
		INSERT INTO transaction.transact_employee_health_record (
			modified_by, created_by, employee_id, 
			mcu_date, periode, status_id, description,
			mcu_url, mcu_followup_url, status_followup_id

		) VALUES ($1, $1, $2, 
		 $3, $4, $5, $6,
		 $7, $8, $9);`
	_, err := tx.Exec(query,
		userId, employeeId,
		value.McuDate, value.Periode, value.StatusId, value.HealthDescription,
		value.McuUrl, value.McuFollowupUrl, value.StatusFollowupId)

	if err != nil {
		return fmt.Errorf("failed to insert transact_employee_health_record: %v", err)
	}

	return nil
}

func (r *repository) UpdateEmployeeHealthRecord(tx *sql.Tx, value HealthRecordModel) error {
	query := `UPDATE transaction.transact_employee_health_record 
		SET modified_date=now(), 
		 modified_by =  $1,
			mcu_date =  $2,
			periode =  $3,
			status_id =  $4,
			description =  $5,
			mcu_url =  $6,
			mcu_followup_url =  $7,
			status_followup_id =  $8
		WHERE id = $9`
	_, err := tx.Exec(query,
		value.ModifiedBy,
		value.McuDate, value.Periode, value.StatusId, value.HealthDescription,
		value.McuUrl, value.McuFollowupUrl, value.StatusFollowupId, value.Id)

	if err != nil {
		return fmt.Errorf("failed to UpdateEmployeeHealthRecord employee health record: %v", err)
	}

	return nil
}

func (r *repository) DeleteEmployeeHealthRecord(tx *sql.Tx, value HealthRecordModel) error {
	query := `UPDATE transaction.transact_employee_health_record 
		SET modified_date=now(), 
			modified_by = $1, 
			is_active=FALSE 
		WHERE id = $2`
	_, err := tx.Exec(query, value.ModifiedBy, value.Id)
	if err != nil {
		return fmt.Errorf("failed to DeleteEmployeeHealthRecord employee health record: %v", err)
	}

	return nil
}
func (r *repository) DeleteEmployeeHealthRecordByEmployeeId(tx *sql.Tx, value HealthRecordModel) error {
	query := `UPDATE transaction.transact_employee_health_record 
		SET modified_date=now(), 
			modified_by = $1, 
			is_active=FALSE 
		WHERE employee_id = $2`
	_, err := tx.Exec(query, value.ModifiedBy, value.EmployeeId)
	if err != nil {
		return fmt.Errorf("failed to DeleteEmployeeHealthRecordByEmployeeId employee health record: %v", err)
	}

	return nil
}
