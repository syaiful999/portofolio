package employee

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func (r *repository) GetEmployeeWorkExperience(ctx context.Context, employeeId uuid.UUID) ([]WorkExperienceModel, error) {

	var items []WorkExperienceModel

	query := `SELECT 
				id, company_name, position, join_date, end_date
			FROM master.master_employee_work_experience
			WHERE is_active = TRUE AND employee_id = $1
			 ORDER BY id`

	rows, err := r.db.Query(query, employeeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var result WorkExperienceModel
		if err := rows.Scan(
			&result.Id, &result.CompanyName, &result.Position, &result.JoinDate, &result.EndDate,
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

func (r *repository) InsertEmployeeMultipleWorkExperiences(tx *sql.Tx, ctx context.Context, userId, employeeId string, values []WorkExperienceModel) error {
	if len(values) == 0 {
		return nil
	}

	query := `
		INSERT INTO "master".master_employee_work_experience (
			created_by, modified_by, employee_id, company_name, position, join_date, end_date
		) VALUES `

	// Build the VALUES clause dynamically
	var params []interface{}
	paramIndex := 1

	for i, value := range values {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			paramIndex, paramIndex+1, paramIndex+2, paramIndex+3, paramIndex+4, paramIndex+5, paramIndex+6)

		params = append(params, userId, userId, employeeId, value.CompanyName, value.Position, value.JoinDate, value.EndDate)
		paramIndex += 7
	}

	_, err := tx.ExecContext(ctx, query, params...)
	if err != nil {
		return fmt.Errorf("failed to insert employee work experiences: %v", err)
	}

	return nil
}

func (r *repository) InsertEmployeeWorkExperience(tx *sql.Tx, userId, employeeId string, value WorkExperienceModel) error {
	query := `
		INSERT INTO master.master_employee_work_experience (
			modified_by, created_by, employee_id, 
			 company_name, position, join_date, end_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7);`

	_, err := tx.Exec(query,
		userId, userId, employeeId,
		value.CompanyName, value.Position, value.JoinDate, value.EndDate,
	)

	if err != nil {
		return fmt.Errorf("failed to insert employee work experience: %v", err)
	}

	return nil
}

func (r *repository) UpdateEmployeeWorkExperience(tx *sql.Tx, value WorkExperienceModel) error {
	query := `UPDATE master.master_employee_work_experience 
		SET modified_date=now(), 
			modified_by =  $1,
			company_name =  $2,
			position =  $3,
			join_date =  $4,
			end_date =  $5
		WHERE id = $6`
	_, err := tx.Exec(query,
		value.ModifiedBy,
		value.CompanyName, value.Position, value.JoinDate, value.EndDate,
		value.Id)

	if err != nil {
		return fmt.Errorf("failed to UpdateEmployeeWorkExperience employee dependents: %v", err)
	}

	return nil
}

func (r *repository) DeleteEmployeeWorkExperience(tx *sql.Tx, value WorkExperienceModel) error {
	query := `UPDATE master.master_employee_work_experience 
		SET modified_date=now(), 
			modified_by = $1, 
			is_active=FALSE 
		WHERE id = $2`
	_, err := tx.Exec(query, value.ModifiedBy, value.Id)
	if err != nil {
		return fmt.Errorf("failed to DeleteEmployeeWorkExperience employee dependents: %v", err)
	}

	return nil
}
