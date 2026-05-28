package employee

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

//

func (r *repository) GetEmployeeCompetency(ctx context.Context, employeeId uuid.UUID) ([]JobCompetencyModel, error) {

	var items []JobCompetencyModel

	query := `SELECT
				id, job_description_id, job_description,
				job_competency_id, competency_name, is_required, total_training
			FROM transaction.v_transact_employee_competency vtec 
			WHERE is_active = TRUE AND employee_id = $1
			 ORDER BY id`

	rows, err := r.db.Query(query, employeeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var result JobCompetencyModel
		if err := rows.Scan(
			&result.Id, &result.JobDescriptionId, &result.JobDescription,
			&result.JobCompetencyId, &result.CompetencyName, &result.IsRequired, &result.TotalTraining,
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

func (r *repository) UpdateEmployeeJobDescription(tx *sql.Tx, value MasterEmployee) error {
	mutation := `UPDATE master.master_employee
				SET modified_by=$1,
					job_description_id=$2,
					modified_date=now()	
				WHERE
					id = $3`
	_, err := tx.Exec(mutation,
		value.ModifiedBy,
		value.JobDescriptionId,
		value.ID,
	)

	return err
}

func (r *repository) InsertEmployeeMultipleCompetencies(tx *sql.Tx, userId, employeeId string, values []JobCompetencyModel) error {
	if len(values) == 0 {
		return nil
	}
	query := `
		INSERT INTO "transaction".transact_employee_competency (
			modified_by, created_by,
			employee_id, job_competency_id, job_description_id, is_required
		) VALUES `
	params := []interface{}{}
	for i, v := range values {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)",
			i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6)

		params = append(params,
			userId,
			userId,
			employeeId,
			v.JobCompetencyId,
			v.JobDescriptionId,
			v.IsRequired,
		)
	}

	// Execute the query without RETURNING clause
	_, err := tx.Exec(query, params...)
	if err != nil {
		return fmt.Errorf("failed to InsertEmployeeMultipleCompetencies transact_employee_competency: %v", err)
	}

	return nil
}

func (r *repository) InsertEmployeeCompetency(tx *sql.Tx, employeeId string, value JobCompetencyModel) error {
	query := `INSERT INTO "transaction".transact_employee_competency 
	(modified_by, created_by,
			employee_id, job_competency_id, job_description_id, is_required)
			VALUES($1, $2, 
			$3, $4, $5, $6, $7)`
	_, err := tx.Exec(query,
		value.ModifiedBy,
		employeeId,
		value.JobCompetencyId.String,
		value.JobDescriptionId.String,
		value.IsRequired,
	)

	if err != nil {
		return fmt.Errorf("failed to InsertEmployeeCompetency employee dependents: %v", err)
	}

	return nil
}

func (r *repository) UpdateEmployeeCompetency(tx *sql.Tx, value JobCompetencyModel) error {
	query := `UPDATE "transaction".transact_employee_competency 
		SET modified_date=now(), 
		 modified_by =  $1,
			job_competency_id =  $2,
			job_description_id =  $3,
			is_required =  $4,
			is_required =  $5
		WHERE id = $6`
	_, err := tx.Exec(query,
		value.ModifiedBy,
		value.JobCompetencyId,
		value.JobDescriptionId,
		value.IsRequired,
		value.IsActive,
		value.Id)

	if err != nil {
		return fmt.Errorf("failed to UpdateEmployeeCompetency employee dependents: %v", err)
	}

	return nil
}

func (r *repository) DeleteEmployeeCompetencyByEmployeeId(tx *sql.Tx, value JobCompetencyModel) error {
	query := `UPDATE "transaction".transact_employee_competency 
		SET modified_date=now(), 
			modified_by = $1, 
			is_active=FALSE 
		WHERE employee_id = $2`
	_, err := tx.Exec(query, value.ModifiedBy, value.EmployeeID)
	if err != nil {
		return fmt.Errorf("failed to DeleteEmployeeCompetencyByEmployeeId employee dependents: %v", err)
	}

	return nil
}
