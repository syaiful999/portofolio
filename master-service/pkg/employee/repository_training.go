package employee

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func (r *repository) GetEmployeeTraining(ctx context.Context, employeeId uuid.UUID) ([]TrainingModel, error) {

	var items []TrainingModel

	query := `SELECT
				id, training_title, training_date, organizer_name, certificate_number,
				training_method, expiry_date, effective_date, cost, certificate_url, competencies_name
			FROM "transaction".v_transact_employee_training
			WHERE is_active = TRUE AND employee_id = $1
			 ORDER BY id`

	rows, err := r.db.Query(query, employeeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var result TrainingModel
		if err := rows.Scan(
			&result.Id, &result.TrainingTitle, &result.TrainingDate, &result.OrganizerName, &result.CertificateNumber, &result.TrainingMethod,
			&result.ExpiryDate, &result.EffectiveDate, &result.Cost, &result.CertificateUrl, &result.CompetenciesName,
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

func (r *repository) GetEmployeeTrainingCompetency(ctx context.Context, trainingId string) ([]TrainingCompetencyModel, error) {

	var items []TrainingCompetencyModel

	query := `SELECT id, training_id, 
		job_competency_id, other_competency, is_active,
		competency_name, description, competency_group, is_required
		FROM "transaction".v_transact_employee_training_competency 
		WHERE is_active = TRUE AND training_id = $1
		ORDER BY id`

	rows, err := r.db.Query(query, trainingId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var result TrainingCompetencyModel
		if err := rows.Scan(
			&result.Id, &result.TrainingId, &result.JobCompetencyId, &result.OtherCompetency, &result.IsActive,
			&result.CompetencyName, &result.Remark, &result.CompetencyGroup, &result.IsRequired,
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

func (r *repository) GetTrainingRecommendation(ctx context.Context) ([]string, error) {

	var items []string

	query := `SELECT training_title FROM master.v_master_training_title`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var result string
		if err := rows.Scan(&result); err != nil {
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

func (r *repository) GetTrainingOrganizerRecommendation(ctx context.Context) ([]string, error) {

	var items []string

	query := `SELECT organizer_name FROM master.v_master_training_organizer`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var result string
		if err := rows.Scan(&result); err != nil {
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

func (r *repository) InsertEmployeeMultipleTrainings(tx *sql.Tx, ctx context.Context, userId, employeeId string, values []TrainingModel) error {
	if len(values) == 0 {
		return nil
	}

	query := `
		INSERT INTO transaction.transact_employee_training (
			modified_by, created_by, employee_id,
			 training_title, training_date, organizer_name, certificate_number, training_method, expiry_date, effective_date, cost, certificate_url
		) VALUES `

	// Build the VALUES clause dynamically
	var params []interface{}
	paramIndex := 1

	for i, value := range values {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			paramIndex, paramIndex+1, paramIndex+2, paramIndex+3, paramIndex+4, paramIndex+5, paramIndex+6,
			paramIndex+7, paramIndex+8, paramIndex+9, paramIndex+10, paramIndex+11)

		params = append(params, userId, userId, employeeId, value.TrainingTitle, value.TrainingDate, value.OrganizerName,
			value.CertificateNumber, value.TrainingMethod, value.ExpiryDate, value.EffectiveDate,
			value.Cost, value.CertificateUrl)
		paramIndex += 12
	}

	_, err := tx.ExecContext(ctx, query, params...)
	if err != nil {
		return fmt.Errorf("failed to insert employee trainings: %v", err)
	}

	return nil
}

func (r *repository) InsertEmployeeTraining(tx *sql.Tx, userId, employeeId string, value TrainingModel) (string, error) {
	query := `
		INSERT INTO transaction.transact_employee_training (
			modified_by, created_by, employee_id, 
			 training_title, training_date, organizer_name, certificate_number, training_method,
			 expiry_date, effective_date, cost, certificate_url)
		VALUES ($1, $1, $2,
		 $3, $4, $5, $6, $7,
		$8, $9, $10, $11)
		RETURNING id;`
	var trainingId string
	err := tx.QueryRow(query,
		userId, employeeId,
		value.TrainingTitle, value.TrainingDate, value.OrganizerName, value.CertificateNumber, value.TrainingMethod,
		value.ExpiryDate, value.EffectiveDate, value.Cost, value.CertificateUrl,
	).Scan(&trainingId)

	if err != nil {
		return "", fmt.Errorf("failed to insert employee training: %v", err)
	}

	return trainingId, nil
}

func (r *repository) UpdateEmployeeTraining(tx *sql.Tx, value TrainingModel) error {
	query := `UPDATE transaction.transact_employee_training 
		SET modified_date=now(), 
		 modified_by =  $1,
			training_title =  $2,
			training_date =  $3,
			organizer_name =  $4,
			certificate_number =  $5,
			training_method = $6,
			expiry_date = $7,
			effective_date = $8,
			cost = $9,
			certificate_url = $10
		WHERE id = $11`
	_, err := tx.Exec(query,
		value.ModifiedBy,
		value.TrainingTitle, value.TrainingDate, value.OrganizerName,
		value.CertificateNumber, value.TrainingMethod, value.ExpiryDate, value.EffectiveDate,
		value.Cost, value.CertificateUrl, value.Id)

	if err != nil {
		return fmt.Errorf("failed to UpdateEmployeeTraining employee dependents: %v", err)
	}

	return nil
}

func (r *repository) DeleteEmployeeTraining(tx *sql.Tx, value TrainingModel) error {
	query := `UPDATE transaction.transact_employee_training 
		SET modified_date=now(), 
			modified_by = $1, 
			is_active=FALSE 
		WHERE id = $2`
	_, err := tx.Exec(query, value.ModifiedBy, value.Id)
	if err != nil {
		return fmt.Errorf("failed to DeleteEmployeeTraining employee dependents: %v", err)
	}

	return nil
}

func (r *repository) DeleteEmployeeTrainingByEmployeeId(tx *sql.Tx, value TrainingModel) error {
	query := `UPDATE transaction.transact_employee_training 
		SET modified_date=now(), 
			modified_by = $1, 
			is_active=FALSE 
		WHERE employee_id = $2`
	_, err := tx.Exec(query, value.ModifiedBy, value.EmployeeId)
	if err != nil {
		return fmt.Errorf("failed to DeleteEmployeeTrainingByEmployeeId employee dependents: %v", err)
	}

	return nil
}

func (r *repository) InsertEmployeeTrainingCompetency(tx *sql.Tx, userId, employeeId string, value TrainingCompetencyModel) error {
	query := `
		INSERT INTO transaction.transact_employee_training_competency (
			modified_by, created_by, employee_id, 
			training_id, job_competency_id, remark, is_required, other_competency)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ;`

	_, err := tx.Exec(query,
		userId, userId, employeeId,
		value.TrainingId, value.JobCompetencyId, value.Remark, value.IsRequired, value.OtherCompetency)

	if err != nil {
		return fmt.Errorf("failed to InsertEmployeeTrainingCompetency employee training competencies: %v", err)
	}

	return nil
}

func (r *repository) UpdateEmployeeTrainingCompetency(tx *sql.Tx, value TrainingCompetencyModel) error {
	query := `UPDATE transaction.transact_employee_training_competency 
		SET modified_date=now(), 
		 modified_by =  $1,
			training_id =  $2,
			job_competency_id =  $3,
			remark =  $4,
			is_required =  $5,
			other_competency = $6,
			is_active = $7
		WHERE id = $8`
	_, err := tx.Exec(query,
		value.ModifiedBy,
		value.TrainingId,
		value.JobCompetencyId,
		value.Remark,
		value.IsRequired,
		value.OtherCompetency,
		value.IsActive,
		value.Id)

	if err != nil {
		return fmt.Errorf("failed to UpdateEmployeeTrainingCompetency employee training competencies: %v", err)
	}

	return nil
}

func (r *repository) DeleteEmployeeTrainingCompetency(tx *sql.Tx, value TrainingCompetencyModel) error {
	query := `UPDATE transaction.transact_employee_training_competency 
		SET modified_date=now(), 
			modified_by = $1, 
			is_active=FALSE 
		WHERE id = $2`
	_, err := tx.Exec(query, value.ModifiedBy, value.Id)
	if err != nil {
		return fmt.Errorf("failed to DeleteEmployeeTrainingCompetency employee training competencies: %v", err)
	}

	return nil
}

func (r *repository) DeleteEmployeeTrainingCompetencyByTrainingId(tx *sql.Tx, value TrainingCompetencyModel) error {
	query := `UPDATE transaction.transact_employee_training_competency 
		SET modified_date=now(), 
			modified_by = $1, 
			is_active=FALSE 
		WHERE training_id = $2`
	_, err := tx.Exec(query, value.ModifiedBy, value.TrainingId)
	if err != nil {
		return fmt.Errorf("failed to DeleteEmployeeTrainingCompetencyByTrainingId employee training competencies: %v", err)
	}

	return nil
}
